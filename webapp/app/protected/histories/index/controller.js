import Ember from 'ember';

export default Ember.Controller.extend({

  modelIsEmpty: Ember.computed.empty('items', 'items'),
  sortableTableConfig: {

    messageConfig: {
      searchLabel: "Search : ",
    },

    customIcons: {
      "sort-asc": "fa fa-caret-up",
      "sort-desc": "fa fa-caret-down",
      "caret": "fa fa-minus",
      "column-visible": "fa fa-minus",
    },

    customClasses: {
      "pageSizeSelectWrapper": "pagination-number"
    }
  },

  data : Ember.computed('model', 'items', function() {
    return this.setData();
  }),

  setData: function() {
    if (!this.get('items')) {
      return;
    }
    var ret = Ember.A([]);
    this.get('items').forEach(function(item) {
      ret.push(Ember.Object.create({
        user: item.get('userFullName'),
        application: item.get('app.displayName'),
        start: window.moment(item.get('startDate')).format('MMMM Do YYYY, h:mm:ss A'),
        end: window.moment(item.get('endDate')).format('MMMM Do YYYY, h:mm:ss A'),
      }));
    });
    this.set('data', ret);
    return ret;
  },

  columns: [
    {
      "propertyName": "user",
      "title": "User",
      "disableFiltering": true,
      "filterWithSelect": false,
    },
    {
      "propertyName": "application",
      "title": "Application",
      "disableFiltering": true,
      "filterWithSelect": false,
    },
    {
      "propertyName": "start",
      "title": "Start Date",
      "disableFiltering": true,
      "filterWithSelect": false,
    },
    {
      "propertyName": "end",
      "title": "End Date",
      "disableFiltering": true,
      "filterWithSelect": false,
    },
  ],
});
