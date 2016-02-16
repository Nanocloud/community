package com.nanocloud.auth.noauthlogged.connection;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Date;
import java.util.concurrent.atomic.AtomicBoolean;
import javax.json.Json;
import javax.json.JsonObject;

import com.nanocloud.auth.noauthlogged.NoAuthLoggedGuacamoleProperties;
import com.nanocloud.auth.noauthlogged.tunnel.ManagedInetGuacamoleSocket;
import com.nanocloud.auth.noauthlogged.tunnel.ManagedSSLGuacamoleSocket;

import org.apache.commons.codec.binary.Base64;
import org.glyptodon.guacamole.GuacamoleException;
import org.glyptodon.guacamole.environment.Environment;
import org.glyptodon.guacamole.environment.LocalEnvironment;
import org.glyptodon.guacamole.net.GuacamoleSocket;
import org.glyptodon.guacamole.net.GuacamoleTunnel;
import org.glyptodon.guacamole.net.SimpleGuacamoleTunnel;
import org.glyptodon.guacamole.net.auth.simple.SimpleConnection;
import org.glyptodon.guacamole.protocol.ConfiguredGuacamoleSocket;
import org.glyptodon.guacamole.protocol.GuacamoleClientInformation;
import org.glyptodon.guacamole.protocol.GuacamoleConfiguration;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class LoggedConnection extends SimpleConnection {

	/**
	 * Logger for this class.
	 */
	private Logger logger = LoggerFactory.getLogger(LoggedConnection.class);

    /**
     * Backing configuration, containing all sensitive information.
     */
    private GuacamoleConfiguration config;
	public LoggedConnection(String name, String identifier, GuacamoleConfiguration config) {
		super(name, identifier, config);

		this.config = config;
	}
	
    /**
     * Task which handles cleanup of a connection associated with some given
     * ActiveConnectionRecord.
     */
    private class ConnectionCleanupTask implements Runnable {

        /**
         * Whether this task has run.
         */
        private final AtomicBoolean hasRun = new AtomicBoolean(false);
		private final String hostname;
		private final Integer port;
		private final String endpoint;
		private final String username;
		private final String password;
		private ActiveConnectionRecord connection;

        public ConnectionCleanupTask(ActiveConnectionRecord connection) throws GuacamoleException {
        	this.connection = connection;

			Environment env = new LocalEnvironment();

			hostname = env.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERURL, "localhost");
			port = env.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPORT, 80);
			endpoint = env.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERENDPOINT, "rpc");
			username = env.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERUSERNAME);
			password = env.getProperty(NoAuthLoggedGuacamoleProperties.NOAUTHLOGGED_SERVERPASSWORD);
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

        @Override
        public void run() {

            // Only run once
            if (!hasRun.compareAndSet(false, true))
                return;

			logger.info("Trying to log history to " + hostname + ":" + port + "/" + endpoint + " with user " + username);

            try {
				String token = login();

				URL myUrl = new URL("http://" + hostname + ":" + port + "/" + endpoint);
    			HttpURLConnection urlConn = (HttpURLConnection)myUrl.openConnection();
    			urlConn.setInstanceFollowRedirects(false);
  			urlConn.setRequestProperty("Authorization", "Bearer " + token);

  			urlConn.setRequestProperty("Content-Type", "application/json");
  			JsonObject params = Json.createObjectBuilder()
  				.add("data", Json.createObjectBuilder()
  						.add("user_id", "94b8e83b-ced3-4259-a3d5-bdc1629272fd")
  						.add("connection_id", this.connection.getConnectionName())
  						.add("start_date", this.connection.getStartDate().toString())
  						.add("end_date", new Date().toString()))
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
					response.append('\r');
				}
				reader.close();

				logger.info("History transmitted to API");

			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}

        }

    }
	
    @Override
    public GuacamoleTunnel connect(GuacamoleClientInformation info)
            throws GuacamoleException {

        Environment env = new LocalEnvironment();
        ActiveConnectionRecord connection = new ActiveConnectionRecord(this.getName());

        // Get guacd connection parameters
        String hostname = env.getProperty(Environment.GUACD_HOSTNAME, "localhost");
        int port = env.getProperty(Environment.GUACD_PORT, 4822);

        GuacamoleSocket socket;

        // Record new active connection
        Runnable cleanupTask = new ConnectionCleanupTask(connection);

        // If guacd requires SSL, use it
        if (env.getProperty(Environment.GUACD_SSL, false))
            socket = new ConfiguredGuacamoleSocket(
                new ManagedSSLGuacamoleSocket(hostname, port, cleanupTask),
                config, info
            );

        // Otherwise, just connect directly via TCP
        else
            socket = new ConfiguredGuacamoleSocket(
                new ManagedInetGuacamoleSocket(hostname, port, cleanupTask),
                config, info
            );

        return new SimpleGuacamoleTunnel(socket);
        
    }
}
