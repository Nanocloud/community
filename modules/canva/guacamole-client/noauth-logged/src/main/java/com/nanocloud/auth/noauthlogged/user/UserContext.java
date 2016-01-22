package com.nanocloud.auth.noauthlogged.user;

import java.io.*;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.Map;

import com.nanocloud.auth.noauthlogged.NoAuthLoggedConfigContentHandler;
import com.nanocloud.auth.noauthlogged.connection.LoggedConnection;
import org.glyptodon.guacamole.GuacamoleException;
import org.glyptodon.guacamole.GuacamoleServerException;
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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xml.sax.InputSource;
import org.xml.sax.SAXException;
import org.xml.sax.XMLReader;
import org.xml.sax.helpers.XMLReaderFactory;

public class UserContext implements org.glyptodon.guacamole.net.auth.UserContext {

	/**
	 * The unique identifier of the root connection group.
	 */
	private static final String ROOT_IDENTIFIER = "ROOT";

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

	/**
	 * Retrieves the configuration file, as defined within guacamole.properties.
	 *
	 * @return The configuration file, as defined within guacamole.properties.
	 * @throws GuacamoleException If an error occurs while reading the
	 *                            property.
	 */
	private File getConfigurationFile() throws GuacamoleException {

		LocalEnvironment environment = new LocalEnvironment();
		// Get config file, defaulting to GUACAMOLE_HOME/noauth-config.xml
		File configFile = environment.getProperty(NOAUTH_CONFIG);
		if (configFile == null)
			configFile = new File(environment.getGuacamoleHome(), DEFAULT_NOAUTH_CONFIG);

		return configFile;

	}

	public synchronized void init() throws GuacamoleException {

		// Get configuration file
		File configFile = getConfigurationFile();
		logger.debug("Reading configuration file: \"{}\"", configFile);

		// Parse document
		try {

			// Set up parser
			NoAuthLoggedConfigContentHandler contentHandler = new NoAuthLoggedConfigContentHandler();

			XMLReader parser = XMLReaderFactory.createXMLReader();
			parser.setContentHandler(contentHandler);

			// Read and parse file
			Reader reader = new BufferedReader(new FileReader(configFile));
			parser.parse(new InputSource(reader));
			reader.close();

			// Init configs
			configs = contentHandler.getConfigs();

		}
		catch (IOException e) {
			throw new GuacamoleServerException("Error reading configuration file.", e);
		}
		catch (SAXException e) {
			throw new GuacamoleServerException("Error parsing XML file.", e);
		}

	}

	public Map<String, GuacamoleConfiguration> getAuthorizedConfigurations() throws GuacamoleException {

		// Check mapping file mod time
		File configFile = getConfigurationFile();
		if (configFile.exists()) {

			// If modified recently, gain exclusive access and recheck
			synchronized (this) {
				if (configFile.exists()) {
					logger.debug("Configuration file \"{}\" has been modified.", configFile);
					init(); // If still not up to date, re-init
				}
			}

		}

		// If no mapping available, report as such
		if (configs == null)
			throw new GuacamoleServerException("Configuration could not be read.");

		return configs;

	}

}