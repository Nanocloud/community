import Ember from 'ember';

export default Ember.Route.extend({
  setupController (controller, model) {
    controller.set('model', model);
    controller.activateRefreshLoop();
  },

  model() {
    return this.store.findAll('machine');
  }
});
