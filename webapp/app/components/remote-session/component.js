import Ember from 'ember';

export default Ember.Component.extend({
  remoteSession: Ember.inject.service('remote-session'),

  guacamole: null,
  connectionName: null,

  getWidth: function() {
    return Ember.$(this.element).parent().width();
  },

  getHeight: function() {
    return Ember.$(this.element).parent().height();
  },

  connect: function() {

    if (Ember.isEmpty(this.get('connectionName'))) {
      return ;
    }

    let width = this.getWidth();
    let height = this.getHeight();

    this.set('guacamole', this.get('remoteSession').getSession(this.get('connectionName'), width, height));
    this.get('guacamole').then((guac) => {

      guac.onfile = function(stream, mimetype, filename) {
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

      this.get('element').appendChild(guac.getDisplay().getElement());

      let mouse = new window.Guacamole.Mouse(guac.getDisplay().getElement());
      let keyboard = new window.Guacamole.Keyboard(document);
      let display = guac.getDisplay();

      window.onresize = function() {
        let width = this.getWidth();
        let height = this.getHeight();

        guac.sendSize(width, height);
      }.bind(this);

      mouse.onmousedown = mouse.onmouseup = mouse.onmousemove = function(mouseState) {
        guac.sendMouseState(mouseState);
      }.bind(this);

      display.oncursor = function(canvas, x, y) {
        display.showCursor(!mouse.setCursor(canvas, x, y));
      };

      keyboard.onkeydown = function (keysym) {
        guac.sendKeyEvent(1, keysym);
      }.bind(this);

      keyboard.onkeyup = function (keysym) {
        guac.sendKeyEvent(0, keysym);
      }.bind(this);


      guac.connect();

      this.get('remoteSession').openedGuacSession[this.get('connectionName')].keyboard = keyboard;

    });
  }.observes('connectionName').on('becameVisible'),

});
