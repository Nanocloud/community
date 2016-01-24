package com.nanocloud.auth.noauthlogged;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;

import com.nanocloud.auth.noauthlogged.user.UserContext;
import org.apache.commons.codec.binary.Base64;
import org.glyptodon.guacamole.GuacamoleException;
import org.glyptodon.guacamole.environment.Environment;
import org.glyptodon.guacamole.environment.LocalEnvironment;
import org.glyptodon.guacamole.net.auth.simple.SimpleAuthenticationProvider;
import org.glyptodon.guacamole.net.auth.AuthenticatedUser;
import org.glyptodon.guacamole.net.auth.Credentials;
import org.glyptodon.guacamole.properties.FileGuacamoleProperty;
import org.glyptodon.guacamole.protocol.GuacamoleConfiguration;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.LoggerFactory;
import org.slf4j.Logger;

import javax.json.Json;
import javax.json.JsonObject;

/**
 * Disable authentication in Guacamole. All users accessing Guacamole are
 * automatically authenticated as "Anonymous" user and are able to use all
 * available GuacamoleConfigurations.
 *
 * GuacamoleConfiguration are read from the XML file defined by `noauth-config`
 * in the Guacamole configuration file (`guacamole.properties`).
 *
 *
 * Example `guacamole.properties`:
 *
 *  noauth-config: /etc/guacamole/noauth-config.xml
 *
 *
 * Example `noauth-config.xml`:
 *
 *  <configs>
 *    <config name="my-rdp-server" protocol="rdp">
 *      <param name="hostname" value="my-rdp-server-hostname" />
 *      <param name="port" value="3389" />
 *    </config>
 *  </configs>
 *
 * @author Laurent Meunier
 */
public class NoAuthLoggedProvider extends SimpleAuthenticationProvider {

    private final String hostname;
    private final Integer port;
    private final String endpoint;
    private final String username;
    private final String password;
    /**
     * Logger for this class.
     */
    private Logger logger = LoggerFactory.getLogger(NoAuthLoggedProvider.class);

    /**
     * The last time the configuration XML was modified, as milliseconds since
     * UNIX epoch.
     */
    private long configTime;

    /**
     * Guacamole server environment.
     */
    private final Environment environment;
    
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
     * Creates a new NoAuthenticationProvider that does not perform any
     * authentication at all. All attempts to access the Guacamole system are
     * presumed to be authorized.
     *
     * @throws GuacamoleException
     *     If a required property is missing, or an error occurs while parsing
     *     a property.
     */
	public NoAuthLoggedProvider() throws GuacamoleException {
		environment = new LocalEnvironment();

        hostname = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERURL, "localhost");
        port = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPORT, 80);
        endpoint = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERENDPOINT, "rpc");
        username = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERUSERNAME);
        password = environment.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPASSWORD);
	}
    
    @Override
    public UserContext getUserContext(AuthenticatedUser authenticatedUser) throws GuacamoleException {

    	Map<String, GuacamoleConfiguration> config = getAuthorizedConfigurations(authenticatedUser.getCredentials());
    	
    	return new UserContext(this, authenticatedUser, config);
    }
    
    @Override
    public String getIdentifier() {
        return "noauthlogged";
    }

    private String login() throws IOException {

        URL myUrl = new URL("http://" + hostname + ":" + port + "/oauth/token");
        HttpURLConnection urlConn = (HttpURLConnection)myUrl.openConnection();
        //urlConn.setInstanceFollowRedirects(false);

        String appKey = "9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae";
        String appSecret = "9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341";
        byte[] auth = Base64.encodeBase64(new String(appKey + ":" + appSecret).getBytes());
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

        URL myUrl = new URL("http://" + hostname + ":" + port + "/api/apps/all");
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

            config.setProtocol("rdp");
            config.setParameter("hostname", connection.getString("Hostname"));
            config.setParameter("port", connection.getString("Port"));
            config.setParameter("username", connection.getString("Username"));
            config.setParameter("password", connection.getString("Password"));
            config.setParameter("security", "nla");
            config.setParameter("ignore-cert", "true");
            if (!connection.has("RemoteApp")) {
                config.setParameter("remote-app", connection.getString("RemoteApp"));
            }

            configs.put(connection.getString("ConnectionName"), config);
        }

        return configs;
    }

    @Override
    public Map<String, GuacamoleConfiguration> getAuthorizedConfigurations(Credentials credentials) throws GuacamoleException {

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
