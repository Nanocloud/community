import Ember from 'ember';
import config from 'nanocloud/config/environment';

/* global Guacamole */

export default Ember.Service.extend({
  session: Ember.inject.service('session'),

  guacamole: null,

  openedGuacSession: {},

  guacToken: function() {
    return Ember.$.post(config.GUACAMOLE_URL + 'api/tokens', {
      access_token: this.get('session.access_token')
    });
  }.property('guacToken'),

  _forgeConnectionString: function(token, connectionName, width, height) {

    // Calculate optimal width/height for display
    var pixel_density = Ember.$(this.get('element')).devicePixelRatio || 1;
    var optimal_dpi = pixel_density * 96;
    var optimal_width = width * pixel_density;
    var optimal_height = height * pixel_density;

    // Build base connect string
    var connectString =
        "token="             + token +
        "&GUAC_DATA_SOURCE=" + "noauthlogged" +
        "&GUAC_ID="          + connectionName +
        "&GUAC_TYPE="        + "c" + // connection
        "&GUAC_WIDTH="       + Math.floor(optimal_width) +
        "&GUAC_HEIGHT="      + Math.floor(optimal_height) +
        "&GUAC_DPI="         + Math.floor(optimal_dpi);

    // Add audio mimetypes to connect string
    connectString += "&GUAC_AUDIO=" + "audio%2Fwav";

    // Add video mimetypes to connect string
    connectString += "&GUAC_VIDEO=" + "video%2Fmp4";

    return connectString;
  },

  getSession: function(name, width, height) {

    return this.get('guacToken').then((token) => {

      let tunnel = new Guacamole.WebSocketTunnel('/guacamole/websocket-tunnel?' + this._forgeConnectionString(token.authToken, name, width, height));
      let guacamole = new Guacamole.Client(
        tunnel
      );

      this.get('openedGuacSession')[name] = { guac : guacamole };

      return guacamole;
    });
  },

  disconnectSession(name) {
      this.get('openedGuacSession')[name].keyboard.onkeydown = null;
      this.get('openedGuacSession')[name].keyboard.onkeyup = null;
      this.get('openedGuacSession')[name].guac.disconnect();
  }
});
