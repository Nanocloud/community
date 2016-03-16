import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  connectionName: null,

  showSidebar: false,
  showSingleTab: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    },
    toggleSingleTab() {
      this.set('connectionName', 'hapticDesktop');

      this.toggleProperty('showSingleTab');
    }
  }
});
