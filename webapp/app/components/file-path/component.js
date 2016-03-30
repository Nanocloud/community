import Ember from 'ember';

export default Ember.Component.extend({
  items: Ember.computed(function() {
    return (this.get('data'));
  })
});
