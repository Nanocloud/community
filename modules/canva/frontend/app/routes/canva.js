import Ember from 'ember';

export default Ember.Route.extend({

  setupController: function(controller, model) {
    controller.set('token', model.token);
    controller.set('connectionName', model.connectionName);
  },
  
  model: function(params) {
    return params;
  }
});
