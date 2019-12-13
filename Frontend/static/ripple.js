/*!
 * ripple.js
 * Copyright (C) 2016-2018 Morgan Bazalgette and Giuseppe Guerra
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// this object contains tiny snippets that were deemed too small to be worth
// their own file.
var singlePageSnippets = {
  "/2fa_gateway" : function() {
    $('#telegram-code')
      .on('input', function() {
        if ($(this).val().length >= 6) {
          $.get("/2fa_gateway/verify", {
            token : $(this).val().trim().substr(0, 8),
          },
          function(resp) {
            switch (resp) {
            case "0":
              $("#telegram-code").closest(".field").addClass("success");
              redir = redir ? redir : "/";
              window.location.href = redir;
              break;
            case "1":
              $("#telegram-code").closest(".field").addClass("error");
              break;
            }
          });
        } else {
          $("#telegram-code").closest(".field").removeClass("error");
        }
      });
  },

  "/leaderboard" : function() {
    page = page === 0 ? 1 : page;

    function loadLeaderboard() {
      var wl = window.location;
      window.history.replaceState(
        '', document.title,
        wl.pathname + "?mode=" + favouriteMode + "&p=" + page +
              (country != "" ? "&country=" + encodeURI(country) : "") +
              wl.hash);
      api("leaderboard", {
        mode : favouriteMode,
        p : page,
        l : 50,
        country : country,
      },
      function(data) {
        var tb = $(".ui.table tbody");
        tb.find("tr").remove();
        if (data.users == null) {
          disableSimplepagButtons(true);
          data.users = [];
        }
        var i = 0;
        data.users.forEach(function(v) {
          tb.append($("<tr />").append(
            $("<td />").text("#" + ((page - 1) * 50 + (++i))),
            $("<td />").html("<a href='/u/" + v.id +
                                   "' title='View profile'><i class='" +
                                   v.country.toLowerCase() + " flag'></i>" +
                                   escapeHTML(v.username) + "</a>"),
            $("<td />").html(
              scoreOrPP(v.chosen_mode.ranked_score, v.chosen_mode.pp)),
            $("<td />").text(v.chosen_mode.accuracy.toFixed(2) + "%"),
            // bonus points if you get the undertale joke
            $("<td />").html(addCommas(v.chosen_mode.playcount) +
                                   " <i title='" + T("Why, LOVE, of course!") +
                                   "'>(lv. " + v.chosen_mode.level.toFixed(0) +
                                   ")</i>")));
        });
        disableSimplepagButtons(data.users.length < 50);
      });
    }
    function scoreOrPP(s, pp) {
      if (pp === 0)
        return "<b>" + addCommas(s) + "</b>";
      return "<b>" + addCommas(pp) + "pp</b> (" + addCommas(s) + ")"
    }

    // country stuff
    $("#country-chooser-modal")
      .click(function() {
        $(".ui.modal").modal("show");
      });
    $(".lb-country")
      .click(function() {
        country = $(this).data("country");
        page = 1;
        $(".ui.modal").modal("hide");
        loadLeaderboard();
      });

    loadLeaderboard();
    setupSimplepag(loadLeaderboard);
    $("#mode-menu .item")
      .click(function(e) {
        e.preventDefault();
        $("#mode-menu .active.item").removeClass("active");
        $(this).addClass("active");
        favouriteMode = $(this).data("mode");
        country = "";
        page = 1;
        loadLeaderboard();
      });
  },

  "/friends" : function() {
    $(".smalltext.button")
      .click(function() {
        var t = $(this);
        var delAdd = t.data("deleted") === "1" ? "add" : "del";
        console.log(delAdd);
        t.addClass("disabled");
        api("friends/" + delAdd, {user : +t.data("userid")}, function(data) {
          t.removeClass("disabled");
          t.data("deleted", data.friend ? "0" : "1");
          t.removeClass("green red blue");
          t.addClass(data.friend ? (data.mutual ? "red" : "green") : "blue");
          t.find(".icon")
            .removeClass("minus plus heart")
            .addClass(data.friend ? (data.mutual ? "heart" : "minus")
              : "plus");
          t.find("span").text(data.friend
            ? (data.mutual ? T("Mutual") : T("Remove"))
            : t("Add"));
        }, true);
      });
  },

  "/team" : function() {
    $("#everyone").click(function() { $(".ui.modal").modal("show"); });
  },

  "/register/verify" : function() {
    var qu = query("u");
    setInterval(function() {
      $.getJSON(hanayoConf.banchoAPI + "/api/v1/verifiedStatus?u=" + qu,
        function(data) {
          if (data.result >= 0) {
            window.location.href = "/register/welcome?u=" + qu;
          }
        })
    }, 5000)
  },

  "/settings" : function() {
    $("input[name='custom_badge.icon']")
      .on("input", function() {
        $("#badge-icon")
          .attr("class", "circular big icon " + escapeHTML($(this).val()));
      });
    $("input[name='custom_badge.name']")
      .on("input", function() {
        $("#badge-name").html(escapeHTML($(this).val()));
      });
    $("input[name='custom_badge.show']")
      .change(function() {
        if ($(this).is(":checked"))
          $("#custom-badge-fields").slideDown();
        else
          $("#custom-badge-fields").slideUp();
      });
    var isDark = $("#dark-site").is(":checked");
    $("form")
      .submit(function(e) {
        e.preventDefault();

        var darkSetting = $("#dark-site").is(":checked")
        if (darkSetting != isDark) {
          var cflags = document.cookie.replace(/(?:(?:^|.*;\s*)cflags\s*\=\s*([^;]*).*$)|^.*$/, "$1");
          cflags = darkSetting ? +cflags | 1 : +cflags & ~1;
          document.cookie = "cflags=" + cflags + ";path=/;max-age=31536000";
        }

        var obj = formToObject($(this));
        var ps = 0;
        $(this)
          .find("input[data-sv]")
          .each(function(_, el) {
            el = $(el);
            if (el.is(":checked")) {
              ps |= el.data("sv");
            }
          });
        obj.play_style = ps;
        var f = $(this);
        api("users/self/settings", obj, function(data) {
          if (darkSetting != isDark) {
            window.location.reload();
            return;
          }
          showMessage("success", "Your new settings have been saved.");
          f.removeClass("loading");
        }, true);
        return false;
      });
  },

  "/settings/userpage" : function() {
    var lastTimeout = null;
    $("textarea[name='data']")
      .on('input', function() {
        if (lastTimeout !== null) {
          clearTimeout(lastTimeout);
        }
        var v = $(this).val();
        lastTimeout = setTimeout(function() {
          $("#userpage-content").addClass("loading");
          $.post(
            "/settings/userpage/parse",
            $("textarea[name='data']").val(), function(data) {
              var e =
                      $("#userpage-content").removeClass("loading").html(data);
              if (typeof twemoji !== "undefined") {
                twemoji.parse(e[0]);
              }
            }, "text");
        }, 800);
      });
    $("form")
      .submit(function(e) {
        e.preventDefault();
        var obj = formToObject($(this));
        var f = $(this);
        api("users/self/userpage", obj, function(data) {
          showMessage("success", "Your userpage has been saved.");
          f.removeClass("loading");
        }, true);
        return false;
      });
  },

  "/donate" : function() {
    var sl = $("#months-slider")[0];
    noUiSlider.create(sl, {
      start : [ 1 ],
      step : 1,
      connect : [ true, false ],
      range : {
        min : [ 1 ],
        max : [ 24 ],
      }
    });
    var rates = {};
    var us = sl.noUiSlider;
    $.getJSON("/donate/rates", function(data) {
      rates = data;
      us.on('update', function() {
        var months = us.get();
        var priceEUR = Math.pow(months * 30 * 0.2, 0.70);
        var priceBTC = priceEUR / rates.EUR;
        var priceUSD = priceBTC * rates.USD;
        $("#cost")
          .html(T("<b>{{ months }}</b> month costs <b>€ {{ eur }}</b>", {
            count : Math.round(+months),
            months : (+months).toFixed(0),
            eur : priceEUR.toFixed(2),
          }) +
                  "<br>" + T("($ {{ usd }} / BTC {{ btc }})", {
              usd : priceUSD.toFixed(2),
              btc : priceBTC.toFixed(10),
            }));
        $("input[name='os0']")
          .attr("value",
            (+months).toFixed(0) + " month" + (months == 1 ? "" : "s"));
        $("#bitcoin-amt").text(priceBTC.toFixed(6));
        $("#paypal-amt").val(priceEUR.toFixed(2));
      });
    });
    $("#username-input").on("input", function() {
      $("#ipn-username").attr("value", "username=" + $(this).val());
    });
  },

  "/settings/avatar" : function() {
    $("#file")
      .change(function(e) {
        var f = e.target.files;
        if (f.length < 1) {
          return;
        }
        var u = window.URL.createObjectURL(f[0]);
        var i = $("#avatar-img")[0];
        i.src = u;
        i.onload = function() { window.URL.revokeObjectURL(this.src); };
      });
  },

  "/beatmaps/rank_request" : function() {
    function updateRankRequestPage(data) {
      $("#queue-info").html(data.submitted + "/" + data.queue_size);

      if (data.submitted_by_user == 0)
        $("#by-you").attr("hidden", "hidden");
      else
        $("#by-you").removeAttr("hidden");

      $("#submitted-by-user").text(data.submitted_by_user);
      $("#max-per-user").text(data.max_per_user);

      var perc = (data.submitted / data.queue_size * 100).toFixed(0);
      $("#progressbar .progress").text(perc + "%");
      $("#progressbar")
        .progress({
          percent : perc,
        });
      if (data.can_submit)
        $("#b-form .input, #b-form .button").removeClass("disabled");
      else
        $("#b-form .input, #b-form .button").addClass("disabled");
    }
    setInterval(function() {
      api("beatmaps/rank_requests/status", {}, updateRankRequestPage);
    }, 10000);
    var re = /^https?:\/\/osu.ppy.sh\/(s|b)\/(\d+)$/gi;
    $("#b-form")
      .submit(function(e) {
        e.preventDefault();
        var v = $("#beatmap").val().trim();
        var reData = re.exec(v);
        re.exec(); // apparently this is always null, idk
        console.log(v, reData);
        if (reData === null) {
          showMessage(
            "error",
            "Please provide a valid link, in the form " +
                    "of either https://osu.ppy.sh/s/&lt;ID&gt; or https://osu.ppy.sh/b/&lt;ID&gt;.");
          $(this).removeClass("loading");
          return false;
        }
        var postData = {};
        if (reData[1] == "s")
          postData.set_id = +reData[2];
        else
          postData.id = +reData[2];
        var t = $(this);
        api("beatmaps/rank_requests", postData,
          function(data) {
            t.removeClass("loading");
            showMessage("success",
              "Beatmap rank request has been submitted.");
            updateRankRequestPage(data);
          },
          function(data) {
            t.removeClass("loading");
            if (data.code == 406)
              showMessage("warning", "That beatmap is already ranked!");
          },
          true);
        return false;
      });
  },

  "/settings/profbackground" : function() {
    $("#colorpicker")
      .minicolors({
        inline : true,
      });
    $("#background-type")
      .change(function() {
        $("[data-type]:not([hidden])").attr("hidden", "hidden");
        $("[data-type=" + $(this).val() + "]").removeAttr("hidden");
      });
    $("#file")
      .change(function(e) {
        var f = e.target.files;
        if (f.length < 1) {
          return;
        }
        var u = window.URL.createObjectURL(f[0]);
        var i = document.createElement("img");
        i.src = u;
        i.onload = function() { window.URL.revokeObjectURL(this.src); };
        $("#image-background").empty().append(i);
      });
  },

  "/dev/tokens" : function() {
    $("#privileges-number")
      .on("input", function() {
        $("#privileges-text").text(privilegesToString($(this).val()));
      });
  }
};

$(document)
  .ready(function() {
    // semantic stuff
    $('.message .close').on('click', closeClosestMessage);
    $('.ui.checkbox').checkbox();
    $('.ui.dropdown').dropdown();
    $('.ui.progress').progress();
    $('.ui.form')
      .submit(function(e) {
        var t = $(this);
        if (t.hasClass("loading") || t.hasClass("disabled")) {
          e.preventDefault();
          return false;
        }
        t.addClass("loading");
        var f = t.attr("id");
        $("[form='" + f + "']").addClass("loading");
      });

    // emojis!
    if (typeof twemoji !== "undefined") {
      $(".twemoji").each(function(k, v) { twemoji.parse(v); });
    }

    // ripple stuff
    var f = singlePageSnippets[window.location.pathname];
    if (typeof f === 'function')
      f();
    if (typeof deferredToPageLoad === "function")
      deferredToPageLoad();

    // setup user search
    $("#user-search")
      .search({
        onSelect : function(val) {
          window.location.href = val.url;
          return false;
        },
        apiSettings : {
          url : "/api/v1/users/lookup?name={query}",
          onResponse : function(resp) {
            var r = {
              results : [],
            };
            $.each(resp.users, function(index, item) {
              r.results.push({
                title : item.username,
                url : "/u/" + item.id,
                image : hanayoConf.avatars + "/" + item.id,
              });
            });
            return r;
          },
        },
      });
    $("#user-search-input")
      .keypress(function(e) {
        if (e.which == 13) {
          window.location.pathname = "/u/" + $(this).val();
        }
      });

    $(document)
      .keydown(function(e) {
        var activeElement = $(document.activeElement);
        var isInput = activeElement.is(":input,[contenteditable]");
        if ((e.which === 83 || e.which === 115) && !isInput) {
          $("#user-search-input").focus();
          e.preventDefault();
        }
        if (e.which === 27 && isInput) {
          activeElement.blur();
        }
      });

    // setup timeago
    $.timeago.settings.allowFuture = true;
    $("time.timeago").timeago();

    $("#language-selector .item")
      .click(function() {
        var lang = $(this).data("lang");
        document.cookie = "language=" + lang + ";path=/;max-age=31536000";
        window.location.reload();
      });
  });

function closeClosestMessage() {
  $(this).closest('.message').fadeOut(300, function() { $(this).remove(); });
};

function showMessage(type, message) {
  var newEl =
      $('<div class="ui ' + type +
        ' message hidden"><i class="close icon"></i>' + T(message) + '</div>');
  newEl.find(".close.icon").click(closeClosestMessage);
  $("#messages-container").append(newEl);
  newEl.slideDown(300);
};

// function for all api calls
function api(endpoint, data, success, failure, post) {
  if (typeof data == "function") {
    success = data;
    data = null;
  }
  if (typeof failure == "boolean") {
    post = failure;
    failure = undefined;
  }

  var errorMessage =
      "An error occurred while contacting the Ripple API. Please report this to a Ripple developer.";

  $.ajax({
    method : (post ? "POST" : "GET"),
    dataType : "json",
    url : hanayoConf.baseAPI + "/api/v1/" + endpoint,
    data : (post ? JSON.stringify(data) : data),
    contentType : (post ? "application/json; charset=utf-8" : ""),
    success : function(data) {
      if (data.code != 200) {
        if ((data.code >= 400 && data.code < 500) &&
            typeof failure == "function") {
          failure(data);
          return;
        }
        console.warn(data);
        showMessage("error", errorMessage);
      }
      success(data);
    },
    error : function(jqXHR, textStatus, errorThrown) {
      if ((jqXHR.status >= 400 && jqXHR.status < 500) &&
          typeof failure == "function") {
        failure(jqXHR.responseJSON);
        return;
      }
      console.warn(jqXHR, textStatus, errorThrown);
      showMessage("error", errorMessage);
    },
  });
};

var modes = {
  0 : "osu! standard",
  1 : "Taiko",
  2 : "Catch the Beat",
  3 : "osu!mania",
};
var modesShort = {
  0 : "std",
  1 : "taiko",
  2 : "ctb",
  3 : "mania",
};

var entityMap = {
  "&" : "&amp;",
  "<" : "&lt;",
  ">" : "&gt;",
  '"' : '&quot;',
  "'" : '&#39;',
  "/" : '&#x2F;',
};
function escapeHTML(str) {
  return String(str).replace(/[&<>"'\/]/g,
    function(s) { return entityMap[s]; });
}

function setupSimplepag(callback) {
  var el = $(".simplepag");
  el.find(".left.floated .item").click(function() {
    if ($(this).hasClass("disabled"))
      return false;
    page--;
    callback();
  });
  el.find(".right.floated .item").click(function() {
    if ($(this).hasClass("disabled"))
      return false;
    page++;
    callback();
  });
}
function disableSimplepagButtons(right) {
  var el = $(".simplepag");

  if (page <= 1)
    el.find(".left.floated .item").addClass("disabled");
  else
    el.find(".left.floated .item").removeClass("disabled");

  if (right)
    el.find(".right.floated .item").addClass("disabled");
  else
    el.find(".right.floated .item").removeClass("disabled");
}

window.URL = window.URL || window.webkitURL;

// thank mr stackoverflow
function addCommas(nStr) {
  nStr += '';
  x = nStr.split('.');
  x1 = x[0];
  x2 = x.length > 1 ? '.' + x[1] : '';
  var rgx = /(\d+)(\d{3})/;
  while (rgx.test(x1)) {
    x1 = x1.replace(rgx, '$1' +
                             ',' +
                             '$2');
  }
  return x1 + x2;
}

// helper functions copied from user.js in old-frontend
function getScoreMods(m, noplus) {
	var r = [];
  // has nc => remove dt
  if ((m & 512) == 512)
    m = m & ~64;
  // has pf => remove sd
  if ((m & 16384) == 16384)
    m = m & ~32;
  modsString.forEach(function(v, idx) {
    var val = 1 << idx;
    if ((m & val) > 0)
      r.push(v);
  });
	if (r.length > 0) {
		return (noplus ? "" : "+ ") + r.join(", ");
	} else {
		return (noplus ? T('None') : '');
	}
}

var modsString = [
  "NF",
	"EZ",
	"NV",
	"HD",
	"HR",
	"SD",
	"DT",
	"RX",
	"HT",
	"NC",
	"FL",
	"AU", // Auto.
	"SO",
	"AP", // Autopilot.
	"PF",
	"K4",
	"K5",
	"K6",
	"K7",
	"K8",
	"K9",
	"RN", // Random
	"LM", // LastMod. Cinema?
	"K9",
	"K0",
	"K1",
	"K3",
	"K2",
];

// time format (seconds -> hh:mm:ss notation)
function timeFormat(t) {
  var h = Math.floor(t / 3600);
  t %= 3600;
  var m = Math.floor(t / 60);
  var s = t % 60;
  var c = "";
  if (h > 0) {
    c += h + ":";
    if (m < 10) {
      c += "0";
    }
    c += m + ":";
  } else {
    c += m + ":";
  }
  if (s < 10) {
    c += "0";
  }
  c += s;
  return c;
}

// http://stackoverflow.com/a/901144/5328069
function query(name, url) {
  if (!url) {
    url = window.location.href;
  }
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
  if (!results)
    return null;
  if (!results[2])
    return '';
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// Useful for forms contacting the Ripple API
function formToObject(form) {
  var inputs = form.find("input, textarea, select");
  var obj = {};
  inputs.each(function(_, el) {
    el = $(el);
    if (el.attr("name") === undefined) {
      return;
    }
    var parts = el.attr("name").split(".");
    var value;
    switch (el.attr("type")) {
    case "checkbox":
      value = el.is(":checked");
      break;
    default:
      switch (el.data("cast")) {
      case "int":
        value = +el.val();
        break;
      default:
        value = el.val();
        break;
      }
      break;
    }
    obj = modifyObjectDynamically(obj, parts, value);
  });
  return obj;
}

// > modifyObjectDynamically({}, ["nice", "meme", "dude"], "lol")
// { nice: { meme: { dude: 'lol' } } }
function modifyObjectDynamically(obj, inds, set) {
  if (inds.length === 1) {
    obj[inds[0]] = set;
  } else if (inds.length > 1) {
    if (typeof obj[inds[0]] !== "object")
      obj[inds[0]] = {};
    obj[inds[0]] = modifyObjectDynamically(obj[inds[0]], inds.slice(1), set);
  }
  return obj;
}

var langWhitelist = [
  "de", "it", "ko", "es", "ru", "pl", "fr", "nl", "sv", "fi", "ro", "ko", "vi"
];
i18next.use(i18nextXHRBackend).init({
  nsSeparator : false,
  keySeparator : false,
  fallbackLng : false,
  lng : hanayoConf.language,
  whitelist : langWhitelist,
  load : "currentOnly",
  backend : {loadPath : "/static/locale/{{lng}}.json"}
});

var i18nLoaded = $.inArray(hanayoConf.language, langWhitelist) === -1;
i18next.on("loaded", function() { i18nLoaded = true });

function T(s, settings) {
  if (typeof settings !== "undefined" &&
      typeof settings.count !== "undefined" &&
      $.inArray(hanayoConf.language, langWhitelist) === -1 &&
      settings.count !== 1)
    s = keyPlurals[s];
  return i18next.t(s, settings);
}

var apiPrivileges = [
  "ReadConfidential", "Write", "ManageBadges", "BetaKeys", "ManageSettings",
  "ViewUserAdvanced", "ManageUser", "ManageRoles", "ManageAPIKeys", "Blog",
  "APIMeta", "Beatmap", "Bancho"
];

function privilegesToString(privs) {
  var privList = [];
  apiPrivileges.forEach(function(value, index) {
    if ((privs & (1 << (index + 1))) != 0)
      privList.push(value);
  });
  return privList.join(", ");
}
