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
    });
  }.observes('connectionName').on('becameVisible'),

});
