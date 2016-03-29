import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  connectionName: null,

  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    }
  }
});
