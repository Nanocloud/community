import DS from 'ember-data';

export default DS.Model.extend({
  email: DS.attr('string'),
  activated: DS.attr('boolean'),
  isadmin: DS.attr('boolean'),
  firstname: DS.attr('string'),
  lastname: DS.attr('string'),
  password: DS.attr('string')
});
