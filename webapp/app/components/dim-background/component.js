import Ember from 'ember';

export default Ember.Component.extend({
  click: function() {
    if (!this.get('preventAction')) {
      this.sendAction();
    }
  },
});
