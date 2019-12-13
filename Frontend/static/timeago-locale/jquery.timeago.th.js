(function (factory) {
  if (typeof define === 'function' && define.amd) {
    define(['jquery'], factory);
  } else if (typeof module === 'object' && typeof module.exports === 'object') {
    factory(require('jquery'));
  } else {
    factory(jQuery);
  }
}(function (jQuery) {
  // Thai
  jQuery.timeago.settings.strings = {
    prefixAgo: null,
    prefixFromNow: null,
    suffixAgo: "ที่แล้ว",
    suffixFromNow: "จากตอนนี้",
    seconds: "น้อยกว่าหนึ่งนาที",
    minute: "หนึ่งนาที",
    minutes: "%d นาที",
    hour: "หนึ่งชั่วโมง",
    hours: " %d ชั่วโมง",
    day: "หนึ่งวัน",
    days: "%d วัน",
    month: "หนึ่งเดือน",
    months: "%d เดือน",
    year: "หนึ่งปี",
    years: "%d ปี",
    wordSeparator: "",
    numbers: []
  };
}));