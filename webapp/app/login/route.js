import Ember from 'ember';

export default Ember.Route.extend({

  queryParams: {
   app: {
      refreshModel: true
    }
  },

  beforeModel(transition) {
    if (transition.queryParams.app) {  
      this.set('directLinkParams', transition.queryParams);
    }
  },

  setupController(controller) {
    controller.reset();
    this._super(...arguments);
  },
});
