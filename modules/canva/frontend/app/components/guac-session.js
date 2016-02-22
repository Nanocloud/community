import Ember from 'ember';

/* global Guacamole */
/* global jQuery */

export default Ember.Component.extend({

  classNames: ['guac-session'],

  guacamole: null,
  guac_token: null,
  connectionname: null,

  _forgeConnectionString: function(token, connectionName) {

    // Calculate optimal width/height for display
    var pixel_density = window.devicePixelRatio || 1;
    var optimal_dpi = pixel_density * 96;
    var optimal_width = window.innerWidth * pixel_density;
    var optimal_height = window.innerHeight * pixel_density;

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

  removeSession: function() {

    if (this.guacamole) {
      this.guacamole.disconnect();
      this.get('element').children().remove();
    }	
  },
  
  startSession: function(token, connectionName) {

    this.removeSession();

    let tunnel = new Guacamole.WebSocketTunnel('wss://localhost/guacamole/websocket-tunnel?' + this._forgeConnectionString(token, connectionName));
    this.guacamole = new Guacamole.Client(
      tunnel
    );
    this.get('element').appendChild(this.guacamole.getDisplay().getElement());
    tunnel.onerror = function(error) {
      console.log("Error " + error);
    };
    tunnel.onstatechange = function(state) {
      console.log("Stage changed " + state);
    };

    jQuery(window).onunload = function() {
      if (this.guacamole) {
	this.guacamole.disconnect();
      }
    }.bind(this);
    
    let mouse = new Guacamole.Mouse(this.guacamole.getDisplay().getElement());
    let keyboard = new Guacamole.Keyboard(document);

    this.guacamole.connect();

    mouse.onmousedown = mouse.onmouseup = mouse.onmousemove = function(mouseState) {
      this.guacamole.sendMouseState(mouseState);
    }.bind(this);

    keyboard.onkeydown = function (keysym) {
      this.guacamole.sendKeyEvent(1, keysym);
    }.bind(this);

    keyboard.onkeyup = function (keysym) {
      this.guacamole.sendKeyEvent(0, keysym);
    }.bind(this);

  },

  triggerSession: function() {
    Ember.run.once(this, 'startSession', this.get('guac_token'), this.get('connectionName'));
  }.observes('guac_token', 'connectionName')
});
