import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller, model) {
    controller.set('model', model.toArray());
  },

  model() {
    return this.store.findAll('application', { reload: true });
  }
});
