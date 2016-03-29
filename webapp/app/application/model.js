import DS from 'ember-data';

export default DS.Model.extend({
  alias: DS.attr('string'),
  collectionName: DS.attr('string'),
  displayName: DS.attr('string'),
  filePath: DS.attr('string'),
  iconContent: DS.attr('string'),
});
