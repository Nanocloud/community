import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller, model) {
    controller.set('drivers', model);
  },
  model() {
   return this.store.findAll('machine-driver');
  }
});
