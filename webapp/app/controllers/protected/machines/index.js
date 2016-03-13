import Ember from 'ember';

export default Ember.Controller.extend({
  machines: Ember.computed('model.@each.isNew', 'model.@each.isDeleted', function() {
    return this.get('model').filterBy('isNew', false).filterBy('isDeleted', false);
  })
});
