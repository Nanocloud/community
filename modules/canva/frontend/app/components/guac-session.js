import Ember from 'ember';

/* global Guacamole */
/* global jQuery */

export default Ember.Component.extend({

  classNames: ['guac-session'],

  guacamole: null,
  guac_token: null,
  connectionname: null,

  modalMessage: "",
  isShowingModal: false,
  actions: {
    toggleErrorModal: function() {
      this.toggleProperty('isShowingModal');
    }
  },

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

    let tunnel = new Guacamole.WebSocketTunnel('/guacamole/websocket-tunnel?' + this._forgeConnectionString(token, connectionName));
    this.guacamole = new Guacamole.Client(
      tunnel
    );
    this.get('element').appendChild(this.guacamole.getDisplay().getElement());

    tunnel.onerror = function(error) {

      if (this.get('isShowingModal') == true)
        return ;

      // Impossible to connect
      if (error.code == 512) {
        this.set('modalMessage', "Cannot connect to remote application");
      } else {
        this.set('modalMessage', error.message);
      }
      this.triggerAction({
        action:'toggleErrorModal',
        target: this,
      });
    }.bind(this);

    tunnel.onstatechange = function(state) {

      if (this.get('isShowingModal') == true)
        return ;

      // If disconnected
      if (state == 2) {
        this.set('modalMessage', "You have been disconnected");
        this.triggerAction({
          action:'toggleErrorModal',
          target: this,
        });
      }
    }.bind(this);

    this.guacamole.onstatechange = function(state) {

      if (state == 3) { // If connected
        this.set('isShowingModal', false);
        return ;
      }

      if (this.get('isShowingModal') == true)
        return ;

      if (state == 1) { // If connecting
        this.set('modalMessage', "Connecting to remote application...");
        this.triggerAction({
          action:'toggleErrorModal',
          target: this,
        });
      }
    }.bind(this);

    this.guacamole.onfile = function(stream, mimetype, filename) {
      let blob_reader = new Guacamole.BlobReader(stream, mimetype);

      blob_reader.onprogress = function() {
        stream.sendAck("Received", Guacamole.Status.Code.SUCCESS);
      }.bind(this);

      blob_reader.onend = function() {
        //Download file in browser
        var element = document.createElement('a');
        element.setAttribute('href', window.URL.createObjectURL(blob_reader.getBlob()));
        element.setAttribute('download', filename);
        element.style.display = 'none';
        document.body.appendChild(element);

        element.click();

        document.body.removeChild(element);
      }.bind(this);

      stream.sendAck("Ready", Guacamole.Status.Code.SUCCESS);
    }.bind(this);

    jQuery(window).onunload = function() {
      if (this.guacamole) {
        this.guacamole.disconnect();
      }
    }.bind(this);

    let mouse = new Guacamole.Mouse(this.guacamole.getDisplay().getElement());
    let keyboard = new Guacamole.Keyboard(document);
    let display = this.guacamole.getDisplay();

    window.onresize = function() {
      let width = window.innerWidth;
      let height = window.innerHeight;
      this.guacamole.sendSize(width, height);
    }.bind(this);

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
    Ember.run.once(this, 'startSession', this.get('guac_token'), this.get('connectionName'));
  }.observes('guac_token', 'connectionName')
});
