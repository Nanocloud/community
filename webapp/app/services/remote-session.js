import Ember from 'ember';
import config from '../config/environment';

/* global Guacamole */

export default Ember.Service.extend({
  session: Ember.inject.service('session'),

  guacToken: null,
  connectionName: null,
  canva: null,
  guacamole: null,

  getGuacToken: function() {
    jQuery.post(config.GUACAMOLE_URL + 'api/tokens', {
      access_token: this.get('session.access_token')
    }, function(post) {
      this.set('guacToken', post.authToken);
    }.bind(this));
  }.on('init'),

  _forgeConnectionString: function(canva, connectionName) {

    // Calculate optimal width/height for display
    var pixel_density = Ember.$(canva.get('element')).devicePixelRatio || 1;
    var optimal_dpi = pixel_density * 96;
    var optimal_width = Ember.$(canva.get('element')).innerWidth() * pixel_density;
    var optimal_height = Ember.$(canva.get('element')).innerHeight() * pixel_density;

    // Build base connect string
    var connectString =
        "token="             + this.get('guacToken') +
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

  connect: function() {

    let canva = this.get('canva');
    let name = this.get('connectionName')

    let tunnel = new Guacamole.WebSocketTunnel('/guacamole/websocket-tunnel?' + this._forgeConnectionString(canva, name));
    this.guacamole = new Guacamole.Client(
      tunnel
    );
    canva.get('element').appendChild(this.guacamole.getDisplay().getElement());

    let mouse = new Guacamole.Mouse(this.guacamole.getDisplay().getElement());
    let keyboard = new Guacamole.Keyboard(document);
    let display = this.guacamole.getDisplay();

    this.guacamole.connect();

    mouse.onmousedown = mouse.onmouseup = mouse.onmousemove = function(mouseState) {
      this.guacamole.sendMouseState(mouseState);
    }.bind(this);

    display.oncursor = function(canvas, x, y) {
      display.showCursor(!mouse.setCursor(canvas, x, y));
    }

    keyboard.onkeydown = function (keysym) {
      this.guacamole.sendKeyEvent(1, keysym);
    }.bind(this);


    keyboard.onkeyup = function (keysym) {
      this.guacamole.sendKeyEvent(0, keysym);
    }.bind(this);
  },

  triggerSession: function() {
    if (this.get('guacToken') && this.get('connectionName') && this.get('canva')) {
      Ember.run.once(this, 'connect');
    }
  }.observes('guacToken', 'connectionName', 'canva')
});
