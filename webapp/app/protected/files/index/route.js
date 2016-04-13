import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller, model) {
    controller.set('items', model);
  },

  model() {
    return this.get('store').query('file', { filename: './' });
  }
});
