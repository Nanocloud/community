import Ember from 'ember';

export default Ember.Controller.extend({
  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    }
  }
});
