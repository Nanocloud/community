import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    }
  }
});
