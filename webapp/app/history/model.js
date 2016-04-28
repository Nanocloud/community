import DS from 'ember-data';

export default DS.Model.extend({
  
    userId: DS.attr('string'),
    connectionId: DS.attr('string'),
    startDate: DS.attr('string'),
    endDate: DS.attr('string'),
    user: DS.belongsTo('user'),
    application: DS.belongsTo('application')
});
