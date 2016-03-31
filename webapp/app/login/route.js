import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller) {
    controller.reset();
    this._super(...arguments);
  },

  redirect() {
    if (this.get('session.isAuthenticated')) {
      this.transitionTo('protected.index');
    }
  }
});
