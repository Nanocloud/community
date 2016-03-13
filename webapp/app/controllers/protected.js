import Ember from 'ember';

export default Ember.Controller.extend({
  showSidebar: false,

  evaluteWindow: function() {

    let testWindow = function() {
      if (Ember.$(window).width() > 991) {
        this.set('showSidebar', true);
      } else {
        this.set('showSidebar', false);
      }
    }.bind(this)

    Ember.$(window).resize(testWindow);

    testWindow();
  }.on("init"),

  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    }
  }
});
