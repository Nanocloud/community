import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  
    userId: DS.attr('string'),
    connectionId: DS.attr('string'),
    startDate: DS.attr('date'),
    endDate: DS.attr('date'),
    user: DS.belongsTo('user'),
    application: DS.belongsTo('application'),
    duration: Ember.computed('startDate', 'endDate', function() {
      var start = window.moment(this.get('startDate'));
      var end = window.moment(this.get('endDate'));
      return end.diff(start);
    })
});
