package com.nanocloud.auth.noauthlogged.user;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.*;

import com.nanocloud.auth.noauthlogged.NoAuthLoggedGuacamoleProperties;
import com.nanocloud.auth.noauthlogged.connection.LoggedConnection;
import org.glyptodon.guacamole.GuacamoleException;
import org.glyptodon.guacamole.environment.LocalEnvironment;
import org.glyptodon.guacamole.form.Form;
import org.glyptodon.guacamole.net.auth.ActiveConnection;
import org.glyptodon.guacamole.net.auth.AuthenticatedUser;
import org.glyptodon.guacamole.net.auth.AuthenticationProvider;
import org.glyptodon.guacamole.net.auth.Connection;
import org.glyptodon.guacamole.net.auth.ConnectionGroup;
import org.glyptodon.guacamole.net.auth.Directory;
import org.glyptodon.guacamole.net.auth.User;
import org.glyptodon.guacamole.net.auth.simple.SimpleConnectionDirectory;
import org.glyptodon.guacamole.net.auth.simple.SimpleConnectionGroup;
import org.glyptodon.guacamole.net.auth.simple.SimpleConnectionGroupDirectory;
import org.glyptodon.guacamole.net.auth.simple.SimpleDirectory;
import org.glyptodon.guacamole.net.auth.simple.SimpleUser;
import org.glyptodon.guacamole.net.auth.simple.SimpleUserDirectory;
import org.glyptodon.guacamole.properties.FileGuacamoleProperty;
import org.glyptodon.guacamole.protocol.GuacamoleConfiguration;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.json.Json;
import javax.json.JsonObject;

public class UserContext implements org.glyptodon.guacamole.net.auth.UserContext {

	private final String hostname;
	private final Integer port;
	private final String endpoint;
	private final String username;
	private final String password;

	/**
	 * The unique identifier of the root connection group.
	 */
	private static final String ROOT_IDENTIFIER = "ROOT";
	private final LocalEnvironment environment;

	/**
	 * Logger for this class.
	 */
	private Logger logger = LoggerFactory.getLogger(UserContext.class);

	/**
	 * Map of all known configurations, indexed by identifier.
	 */
	private Map<String, GuacamoleConfiguration> configs;

	private User self;
	private AuthenticationProvider authProvider;
	private Directory<User> userDirectory;
	private Directory<Connection> connectionDirectory;
	private Directory<ConnectionGroup> connectionGroupDirectory;
	private ConnectionGroup rootGroup;

	public UserContext(AuthenticationProvider authProvider,
					   AuthenticatedUser authenticatedUser, Map<String, GuacamoleConfiguration> configs) throws GuacamoleException {

		environment = new LocalEnvironment();

		hostname = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERURL, "localhost");
		port = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPORT, 80);
		endpoint = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERENDPOINT, "rpc");
		username = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERUSERNAME);
		password = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPASSWORD);

		// Return as unauthorized if not authorized to retrieve configs
		if (configs == null)
			throw new GuacamoleException("No configuration file");

		Collection<String> connectionIdentifiers = new ArrayList<String>(configs.size());
		Collection<String> connectionGroupIdentifiers = Collections.singleton(ROOT_IDENTIFIER);

		// Produce collection of connections from given configs
		Collection<Connection> connections = new ArrayList<Connection>(configs.size());
		for (Map.Entry<String, GuacamoleConfiguration> configEntry : configs.entrySet()) {

			// Get connection identifier and configuration
			String identifier = configEntry.getKey();
			GuacamoleConfiguration config = configEntry.getValue();

			// Add as simple connection
			Connection connection = new LoggedConnection(identifier, identifier, config);
			connection.setParentIdentifier(ROOT_IDENTIFIER);
			connections.add(connection);

			// Add identifier to overall set of identifiers
			connectionIdentifiers.add(identifier);

		}

		// Add root group that contains only the given configurations
		this.rootGroup = new SimpleConnectionGroup(
				ROOT_IDENTIFIER, ROOT_IDENTIFIER,
				connectionIdentifiers, Collections.<String>emptyList()
		);

		// Build new user from credentials
		this.self = new SimpleUser(authenticatedUser.getIdentifier(), connectionIdentifiers,
				connectionGroupIdentifiers);

		// Create directories for new user
		this.userDirectory = new SimpleUserDirectory(self);
		this.connectionDirectory = new SimpleConnectionDirectory(connections);
		this.connectionGroupDirectory = new SimpleConnectionGroupDirectory(Collections.singleton(this.rootGroup));

		// Associate provided AuthenticationProvider
		this.authProvider = authProvider;
	}

	@Override
	public User self() {
		return self;
	}

	@Override
	public AuthenticationProvider getAuthenticationProvider() {
		return authProvider;
	}

	@Override
	public Directory<User> getUserDirectory() throws GuacamoleException {
		return userDirectory;
	}

	@Override
	public Directory<Connection> getConnectionDirectory() throws GuacamoleException {

		System.out.println("UserContext : getConnectionDirectory");
		configs = this.getAuthorizedConfigurations();
		Collection<String> connectionIdentifiers = new ArrayList<String>(configs.size());

		// Produce collection of connections from given configs
		Collection<Connection> connections = new ArrayList<Connection>(configs.size());
		for (Map.Entry<String, GuacamoleConfiguration> configEntry : configs.entrySet()) {

			// Get connection identifier and configuration
			String identifier = configEntry.getKey();
			GuacamoleConfiguration config = configEntry.getValue();

			// Add as simple connection
			Connection connection = new LoggedConnection(identifier, identifier, config);
			connection.setParentIdentifier(ROOT_IDENTIFIER);
			connections.add(connection);
			System.out.println("UserContext : " + connection.getName());

			// Add identifier to overall set of identifiers
			connectionIdentifiers.add(identifier);
		}

		this.connectionDirectory = new SimpleConnectionDirectory(connections);

		return connectionDirectory;
	}

	@Override
	public Directory<ConnectionGroup> getConnectionGroupDirectory() throws GuacamoleException {
		return connectionGroupDirectory;
	}

	@Override
	public Directory<ActiveConnection> getActiveConnectionDirectory() throws GuacamoleException {
		return new SimpleDirectory<ActiveConnection>();
	}

	@Override
	public ConnectionGroup getRootConnectionGroup() throws GuacamoleException {
		return rootGroup;
	}

	@Override
	public Collection<Form> getUserAttributes() {
		return Collections.<Form>emptyList();
	}

	@Override
	public Collection<Form> getConnectionAttributes() {
		return Collections.<Form>emptyList();
	}

	@Override
	public Collection<Form> getConnectionGroupAttributes() {
		return Collections.<Form>emptyList();
	}

	/**
	 * The XML file to read the configuration from.
	 */
	public static final FileGuacamoleProperty NOAUTH_CONFIG = new FileGuacamoleProperty() {

		@Override
		public String getName() {
			return "noauthlogged-config";
		}

	};

	/**
	 * The default filename to use for the configuration, if not defined within
	 * guacamole.properties.
	 */
	public static final String DEFAULT_NOAUTH_CONFIG = "noauth-config.xml";

	private String login() throws IOException {

		URL myUrl = new URL("http://" + hostname + ":" + port + "/oauth/token");
		HttpURLConnection urlConn = (HttpURLConnection)myUrl.openConnection();
		//urlConn.setInstanceFollowRedirects(false);

		String appKey = "9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae";
		String appSecret = "9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341";
		byte[] auth = org.apache.commons.codec.binary.Base64.encodeBase64(new String(appKey + ":" + appSecret).getBytes());
		urlConn.setRequestProperty("Authorization", "Basic OTQwNWZiNmIwZTU5ZDI5OTdlM2M3NzdhMjJkOGYwZTYxN2E5ZjViMzZiNjU2NWM3NTc5ZTViZTZkZWI4ZjdhZTo5MDUwZDY3YzJiZTA5NDNmMmM2MzUwNzA1MmRkZWRiM2FlMzRhMzBlMzliYmJiZGFiMjQxYzkzZjhiNWNmMzQx");
		urlConn.setRequestProperty("Content-Type", "application/json");

		JsonObject params = Json.createObjectBuilder()
				.add("username", username)
				.add("password", password)
				.add("grant_type", "password")
				.build();

		urlConn.setUseCaches(false);
		urlConn.setDoOutput(true);
		// Send request (for some reason we actually need to wait for response)
		DataOutputStream writer = new DataOutputStream(urlConn.getOutputStream());
		writer.writeBytes(params.toString());
		writer.close();

		urlConn.connect();
		urlConn.getOutputStream().close();

		// Get Response
		InputStream input = urlConn.getInputStream();
		BufferedReader reader = new BufferedReader(new InputStreamReader(input));
		StringBuilder response = new StringBuilder();

		String line;
		while ((line = reader.readLine()) != null) {
			response.append(line);
		}
		reader.close();

		JSONObject json = null;
		try {
			json = new JSONObject(response.toString());
		} catch (JSONException e) {
			e.printStackTrace();
		}
		String token = null;
		try {
			token = json.getString("access_token");
		} catch (JSONException e) {
			e.printStackTrace();
		}

		return token;
	}

	private Map<String, GuacamoleConfiguration> askForConnections() throws IOException, JSONException {

		Map<String, GuacamoleConfiguration> configs = new HashMap<String, GuacamoleConfiguration>();

		String token = login();

		URL myUrl = new URL("http://" + hostname + ":" + port + "/api/apps/connections");
		HttpURLConnection urlConn = (HttpURLConnection)myUrl.openConnection();

		urlConn.setInstanceFollowRedirects(false);
		urlConn.setRequestProperty("Authorization", "Bearer " + token);
		urlConn.setUseCaches(false);

		urlConn.connect();

		// Get Response
		InputStream input = urlConn.getInputStream();

		BufferedReader reader = new BufferedReader(new InputStreamReader(input));
		StringBuilder response = new StringBuilder();

		String line;
		while ((line = reader.readLine()) != null) {
			response.append(line);
		}
		reader.close();

		System.out.println(response.toString());
		JSONArray appList =  new JSONArray(response.toString());
		for (int i = 0; i < appList.length(); i++) {
			JSONObject connection = appList.getJSONObject(i);
			GuacamoleConfiguration config = new GuacamoleConfiguration();

			config.setProtocol(connection.getString("protocol"));
			config.setParameter("hostname", connection.getString("hostname"));
			config.setParameter("port", connection.getString("port"));
			config.setParameter("username", connection.getString("username"));
			config.setParameter("password", connection.getString("password"));
			config.setParameter("security", "nla");
			config.setParameter("ignore-cert", "true");
			if (connection.has("remote_app")) {
				config.setParameter("remote-app", connection.getString("remote_app"));
			}

			configs.put(connection.getString("app_name"), config);
		}

		return configs;
	}

	public Map<String, GuacamoleConfiguration> getAuthorizedConfigurations() throws GuacamoleException {

		Map<String, GuacamoleConfiguration> configs = null;

		logger.info("Fetch application list from server");

		try {
			configs = askForConnections();
		} catch (IOException e) {
			e.printStackTrace();
			return null;
		} catch (JSONException e) {
			e.printStackTrace();
		}

		logger.info("Application list fetched");

		return configs;

	}

}