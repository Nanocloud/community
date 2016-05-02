import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller, model) {
    controller.set('types', model);
  },

  model() {
    return this.store.query(
      'machine-type',
      {
        driver: 'aws'
      }
    );
  },

  actions: {
    willTransition: function() {
      this.controller.reset();
    }
  }
});
