import Ember from 'ember';

export default Ember.Controller.extend({
  connectionName: null,

  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    },

    logout() {
      this.get('session').invalidate();
    }
  }
});
