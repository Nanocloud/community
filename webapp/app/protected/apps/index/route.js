import Ember from 'ember';

export default Ember.Route.extend({

  configuration: Ember.inject.service('configuration'),
  model() {
    this.get('configuration').loadData();
    return this.store.findAll('application', { reload: true });
  },

  actions: {
    refreshModel() {
      this.refresh();
    }
  }
});
