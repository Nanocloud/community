import Ember from 'ember';

export default Ember.Route.extend({
  session: Ember.inject.service('session'),

  beforeModel() {
    if (this.get('session.isAuthenticated') === false) {
      this.transitionTo('login');
    }
  },

  actions: {
    invalidateSession: function() {
      this.get('session').invalidate();
      this.transitionTo('login');
    }
  }
});
