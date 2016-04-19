import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return this.store.findAll('history', { reload: true });
  }
});
