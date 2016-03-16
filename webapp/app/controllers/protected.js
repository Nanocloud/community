import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),
  remoteSession: Ember.inject.service('remote-session'),

  guacamole: null,

  showSidebar: false,
  showSingleTab: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    },
    toggleSingleTab() {
      this.set('guacamole', this.get('remoteSession').getSession('hapticDesktop', 800, 600));

      this.toggleProperty('showSingleTab');
    }
  }
});
