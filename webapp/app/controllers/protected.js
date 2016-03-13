import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),
  remoteSession: Ember.inject.service('remote-session'),

  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    }
  }
});
