import Ember from 'ember';

export default Ember.Route.extend({

  queryParams: {
   app: {
      refreshModel: true
    }
  },

  setupController(controller) {
    controller.reset();
    this._super(...arguments);
  },
});
