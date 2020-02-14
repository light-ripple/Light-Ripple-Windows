-- --------------------------------------------------------
-- Host:                         Localhost
-- Server version:               Bedrock Linux 0.7.12 Poki
-- Server OS:                    Linux
-- osu!DB structure For Windows(localhost) by Light Ripple
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Dumping database structure for Localhost
-- Dumping database structure for Localhost (Win)
CREATE DATABASE IF NOT EXISTS `ripple` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `ripple`;

-- Dumping structure 2fa
CREATE TABLE IF NOT EXISTS `2fa` (
  `userid` int(11) NOT NULL,
  `ip` int(11) NOT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data 2fa: ~0 rows (approximately)
/*!40000 ALTER TABLE `2fa` DISABLE KEYS */;
/*!40000 ALTER TABLE `2fa` ENABLE KEYS */;

-- Dumping structure 2fa_telegram
CREATE TABLE IF NOT EXISTS `2fa_telegram` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Dumping data 2fa_telegram: ~0 rows (approximately)
/*!40000 ALTER TABLE `2fa_telegram` DISABLE KEYS */;
/*!40000 ALTER TABLE `2fa_telegram` ENABLE KEYS */;

-- Dumping structure 2fa_totp
CREATE TABLE IF NOT EXISTS `2fa_totp` (
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `userid` int(11) NOT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data 2fa_totp: ~0 rows (approximately)
/*!40000 ALTER TABLE `2fa_totp` DISABLE KEYS */;
/*!40000 ALTER TABLE `2fa_totp` ENABLE KEYS */;

-- Dumping structure achievements
CREATE TABLE IF NOT EXISTS `achievements` (
  `id` int(11) NOT NULL,
  `name` text CHARACTER SET latin1 NOT NULL,
  `description` text CHARACTER SET latin1 NOT NULL,
  `icon` text CHARACTER SET latin1 NOT NULL,
  `version` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data achievements: ~97 rows (approximately)
/*!40000 ALTER TABLE `achievements` DISABLE KEYS */;
INSERT INTO `achievements` (`id`, `name`, `description`, `icon`, `version`) VALUES
	(1, '500 Combo (osu!std)', '500 big ones! You\'re moving up in the world!', 'osu-combo-500', 1),
	(2, '750 Combo (osu!std)', '750 big ones! You\'re moving up in the world!', 'osu-combo-750', 1),
	(3, '1000 Combo (osu!std)', '1000 big ones! You\'re moving up in the world!', 'osu-combo-1000', 1),
	(4, '2000 Combo (osu!std)', '2000 big ones! You\'re moving up in the world!', 'osu-combo-2000', 1),
	(5, '500 Combo (osu!taiko)', '500 big ones! You\'re moving up in the world!', 'osu-combo-500', 1),
	(6, '750 Combo (osu!taiko)', '750 big ones! You\'re moving up in the world!', 'osu-combo-750', 1),
	(7, '1000 Combo (osu!taiko)', '1000 big ones! You\'re moving up in the world!', 'osu-combo-1000', 1),
	(8, '2000 Combo (osu!taiko)', '2000 big ones! You\'re moving up in the world!', 'osu-combo-2000', 1),
	(9, '500 Combo (osu!ctb)', '500 big ones! You\'re moving up in the world!', 'osu-combo-500', 1),
	(10, '750 Combo (osu!ctb)', '750 big ones! You\'re moving up in the world!', 'osu-combo-750', 1),
	(11, '1000 Combo (osu!ctb)', '1000 big ones! You\'re moving up in the world!', 'osu-combo-1000', 1),
	(12, '2000 Combo (osu!ctb)', '2000 big ones! You\'re moving up in the world!', 'osu-combo-2000', 1),
	(13, '500 Combo (osu!mania)', '500 big ones! You\'re moving up in the world!', 'osu-combo-500', 1),
	(14, '750 Combo (osu!mania)', '750 big ones! You\'re moving up in the world!', 'osu-combo-750', 1),
	(15, '1000 Combo (osu!mania)', '1000 big ones! You\'re moving up in the world!', 'osu-combo-1000', 1),
	(16, '2000 Combo (osu!mania)', '2000 big ones! You\'re moving up in the world!', 'osu-combo-2000', 1),
	(17, 'Rising Star', 'Can\'t go forward without the first steps.', 'osu-skill-pass-1', 2),
	(18, 'My First Don', 'Can\'t go forward without the first steps.', 'taiko-skill-pass-1', 2),
	(19, 'A Slice Of Life', 'Can\'t go forward without the first steps.', 'fruits-skill-pass-1', 2),
	(20, 'First Steps', 'Can\'t go forward without the first steps.', 'mania-skill-pass-1', 2),
	(21, 'Constellation Prize', 'Definitely not a consolation prize. Now things start getting hard!', 'osu-skill-pass-2', 2),
	(22, 'Katsu Katsu Katsu', 'Definitely not a consolation prize. Now things start getting hard!', 'taiko-skill-pass-2', 2),
	(23, 'Dashing Ever Forward', 'Definitely not a consolation prize. Now things start getting hard!', 'fruits-skill-pass-2', 2),
	(24, 'No Normal Player', 'Definitely not a consolation prize. Now things start getting hard!', 'mania-skill-pass-2', 2),
	(25, 'Building Confidence', 'Oh, you\'ve SO got this.', 'osu-skill-pass-3', 2),
	(26, 'Not Even Trying', 'Oh, you\'ve SO got this.', 'taiko-skill-pass-3', 2),
	(27, 'Zesty Disposition', 'Oh, you\'ve SO got this.', 'fruits-skill-pass-3', 2),
	(28, 'Impulse Drive', 'Oh, you\'ve SO got this.', 'mania-skill-pass-3', 2),
	(29, 'Insanity Approaches', 'You\'re not twitching, you\'re just ready.', 'osu-skill-pass-4', 2),
	(30, 'Face Your Demons', 'You\'re not twitching, you\'re just ready.', 'taiko-skill-pass-4', 2),
	(31, 'Hyperdash ON!', 'You\'re not twitching, you\'re just ready.', 'fruits-skill-pass-4', 2),
	(32, 'Hyperspeed', 'You\'re not twitching, you\'re just ready.', 'mania-skill-pass-4', 2),
	(33, 'These Clarion Skies', 'Everything seems so clear now.', 'osu-skill-pass-5', 2),
	(34, 'The Demon Within', 'Everything seems so clear now.', 'taiko-skill-pass-5', 2),
	(35, 'It\'s Raining Fruit', 'Everything seems so clear now.', 'fruits-skill-pass-5', 2),
	(36, 'Ever Onwards', 'Everything seems so clear now.', 'mania-skill-pass-5', 2),
	(37, 'Above and Beyond', 'A cut above the rest.', 'osu-skill-pass-6', 2),
	(38, 'Drumbreaker', 'A cut above the rest.', 'taiko-skill-pass-6', 2),
	(39, 'Fruit Ninja', 'A cut above the rest.', 'fruits-skill-pass-6', 2),
	(40, 'Another Surpassed', 'A cut above the rest.', 'mania-skill-pass-6', 2),
	(41, 'Supremacy', 'All marvel before your prowess.', 'osu-skill-pass-7', 2),
	(42, 'The Godfather', 'All marvel before your prowess.', 'taiko-skill-pass-7', 2),
	(43, 'Dreamcatcher', 'All marvel before your prowess.', 'fruits-skill-pass-7', 2),
	(44, 'Extra Credit', 'All marvel before your prowess.', 'mania-skill-pass-7', 2),
	(45, 'Absolution', 'My god, you\'re full of stars!', 'osu-skill-pass-8', 2),
	(46, 'Rhythm Incarnate', 'My god, you\'re full of stars!', 'taiko-skill-pass-8', 2),
	(47, 'Lord of the Catch', 'My god, you\'re full of stars!', 'fruits-skill-pass-8', 2),
	(48, 'Maniac', 'My god, you\'re full of stars!', 'mania-skill-pass-8', 2),
	(49, 'Totality', 'All the notes. Every single one.', 'osu-skill-fc-1', 3),
	(50, 'Keeping Time', 'All the notes. Every single one.', 'taiko-skill-fc-1', 3),
	(51, 'Sweet And Sour', 'All the notes. Every single one.', 'fruits-skill-fc-1', 3),
	(52, 'Keystruck', 'All the notes. Every single one.', 'mania-skill-fc-1', 3),
	(53, 'Business As Usual', 'Two to go, please.', 'osu-skill-fc-2', 3),
	(54, 'To Your Own Beat', 'Two to go, please.', 'taiko-skill-fc-2', 3),
	(55, 'Reaching The Core', 'Two to go, please.', 'fruits-skill-fc-2', 3),
	(56, 'Keying In', 'Two to go, please.', 'mania-skill-fc-2', 3),
	(57, 'Building Steam', 'Hey, this isn\'t so bad.', 'osu-skill-fc-3', 3),
	(58, 'Big Drums', 'Hey, this isn\'t so bad.', 'taiko-skill-fc-3', 3),
	(59, 'Clean Platter', 'Hey, this isn\'t so bad.', 'fruits-skill-fc-3', 3),
	(60, 'Hyperflow', 'Hey, this isn\'t so bad.', 'mania-skill-fc-3', 3),
	(61, 'Moving Forward', 'Bet you feel good about that.', 'osu-skill-fc-4', 3),
	(62, 'Adversity Overcome', 'Bet you feel good about that.', 'taiko-skill-fc-4', 3),
	(63, 'Between The Rain', 'Bet you feel good about that.', 'fruits-skill-fc-4', 3),
	(64, 'Breakthrough', 'Bet you feel good about that.', 'mania-skill-fc-4', 3),
	(65, 'Paradigm Shift', 'Surprisingly difficult.', 'osu-skill-fc-5', 3),
	(66, 'Demonslayer', 'Surprisingly difficult.', 'taiko-skill-fc-5', 3),
	(67, 'Addicted', 'Surprisingly difficult.', 'fruits-skill-fc-5', 3),
	(68, 'Everything Extra', 'Surprisingly difficult.', 'mania-skill-fc-5', 3),
	(69, 'Anguish Quelled', 'Don\'t choke.', 'osu-skill-fc-6', 3),
	(70, 'Rhythm\'s Call', 'Don\'t choke.', 'taiko-skill-fc-6', 3),
	(71, 'Quickening', 'Don\'t choke.', 'fruits-skill-fc-6', 3),
	(72, 'Level Breaker', 'Don\'t choke.', 'mania-skill-fc-6', 3),
	(73, 'Never Give Up', 'Excellence is its own reward.', 'osu-skill-fc-7', 3),
	(74, 'Time Everlasting', 'Excellence is its own reward.', 'taiko-skill-fc-7', 3),
	(75, 'Supersonic', 'Excellence is its own reward.', 'fruits-skill-fc-7', 3),
	(76, 'Step Up', 'Excellence is its own reward.', 'mania-skill-fc-7', 3),
	(77, 'Aberration', 'They said it couldn\'t be done. They were wrong.', 'osu-skill-fc-8', 3),
	(78, 'The Drummer\'s Throne', 'They said it couldn\'t be done. They were wrong.', 'taiko-skill-fc-8', 3),
	(79, 'Dashing Scarlet', 'They said it couldn\'t be done. They were wrong.', 'fruits-skill-fc-8', 3),
	(80, 'Behind The Veil', 'They said it couldn\'t be done. They were wrong.', 'mania-skill-fc-8', 3),
	(81, 'Finality', 'High stakes, no regrets.', 'all-intro-suddendeath', 4),
	(82, 'Perfectionist', 'Accept nothing but the best.', 'all-intro-perfect', 4),
	(83, 'Rock Around The Clock', 'You can\'t stop the rock.', 'all-intro-hardrock', 4),
	(84, 'Time And A Half', 'Having a right ol\' time. One and a half of them, almost.', 'all-intro-doubletime', 4),
	(85, 'Sweet Rave Party', 'Founded in the fine tradition of changing things that were just fine as they were.', 'all-intro-nightcore', 4),
	(86, 'Blindsight', 'I can see just perfectly.', 'all-intro-hidden', 4),
	(87, 'Are You Afraid Of The Dark?', 'Harder than it looks, probably because it\'s hard to look.', 'all-intro-flashlight', 4),
	(88, 'Dial It Right Back', 'Sometimes you just want to take it easy.', 'all-intro-easy', 4),
	(89, 'Risk Averse', 'Safety nets are fun!', 'all-intro-nofail', 4),
	(90, 'Slowboat', 'You got there. Eventually.', 'all-intro-halftime', 4),
	(91, 'Burned Out', 'One cannot always spin to win.', 'all-intro-spunout', 4),
	(92, '5,000 Plays', 'There\'s a lot more where that came from.', 'osu-plays-5000', 5),
	(93, '15,000 Plays', 'Must.. click.. circles..', 'osu-plays-15000', 5),
	(94, '25,000 Plays', 'There\'s no going back.', 'osu-plays-25000', 5),
	(95, '50,000 Plays', 'You\'re here forever.', 'osu-plays-50000', 5),
	(96, 'A meganekko approaches', 'Congratulations, you met Maria!', 'mania-secret-meganekko', 6),
	(97, 'Don\'t let the bunny distract you!', 'The order was indeed, not a rabbit.', 'all-secret-bunny', 6);
/*!40000 ALTER TABLE `achievements` ENABLE KEYS */;

-- Dumping structure badges
CREATE TABLE IF NOT EXISTS `badges` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(21485) NOT NULL,
  `icon` varchar(32) CHARACTER SET latin1 NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1005 DEFAULT CHARSET=utf8;

-- Dumping data badges: ~13 rows (approximately)
/*!40000 ALTER TABLE `badges` DISABLE KEYS */;
INSERT INTO `badges` (`id`, `name`, `icon`) VALUES
	(0, '', ''),
	(2, 'Developers', 'teal blind'),
	(3, 'Bug Hunter', 'bug icon'),
	(4, 'Community Manager', 'user secret'),
	(5, 'Beatmap Nominators', 'chart line'),
	(10, 'SUSPICIOUS - WAITING FOR CHECK', 'red window close outline'),
	(30, 'Chat Moderators', 'envelope outline'),
	(1000, 'Thumbnail Maker', 'fa-thumbs-o-up'),
	(1001, 'Marathon Runner', 'yellow hourglass outline'),
	(1002, 'Donor', 'yellow heart'),
	(1003, 'Vanilla God (Certified by Aoba)', 'yellow fa-check');
/*!40000 ALTER TABLE `badges` ENABLE KEYS */;

-- Dumping structure bancho_channels
CREATE TABLE IF NOT EXISTS `bancho_channels` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `description` varchar(127) NOT NULL,
  `public_read` tinyint(4) NOT NULL,
  `public_write` tinyint(4) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `temp` tinyint(1) NOT NULL DEFAULT '0',
  `hidden` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- Dumping data bancho_channels: ~7 rows (approximately)
/*!40000 ALTER TABLE `bancho_channels` DISABLE KEYS */;
INSERT INTO `bancho_channels` (`id`, `name`, `description`, `public_read`, `public_write`, `status`, `temp`, `hidden`) VALUES
	(1, '#osu', 'Ainu global chat', 1, 1, 1, 0, 0),
	(2, '#announce', 'Announce channel', 1, 0, 1, 0, 0),
	(3, '#english', 'English community', 1, 1, 1, 0, 0),
	(4, '#admin', 'Are you admin?', 1, 1, 1, 0, 1),
	(5, '#lobby', 'This is the lobby where you find games to play with others!', 1, 1, 1, 0, 1),
	(6, '#ranked', 'Rank requests maps will be posted here! (If it\'s ranked.)', 1, 0, 1, 0, 0);
/*!40000 ALTER TABLE `bancho_channels` ENABLE KEYS */;

-- Dumping structure bancho_messages
CREATE TABLE IF NOT EXISTS `bancho_messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `msg_from_userid` int(11) NOT NULL,
  `msg_from_username` varchar(30) CHARACTER SET latin1 NOT NULL,
  `msg_to` varchar(32) CHARACTER SET latin1 NOT NULL,
  `msg` varchar(127) CHARACTER SET latin1 NOT NULL,
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data bancho_messages: ~0 rows (approximately)
/*!40000 ALTER TABLE `bancho_messages` DISABLE KEYS */;
/*!40000 ALTER TABLE `bancho_messages` ENABLE KEYS */;

-- Dumping structure bancho_private_messages
CREATE TABLE IF NOT EXISTS `bancho_private_messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `msg_from_userid` int(11) NOT NULL,
  `msg_from_username` varchar(30) CHARACTER SET latin1 NOT NULL,
  `msg_to` varchar(32) CHARACTER SET latin1 NOT NULL,
  `msg` varchar(127) CHARACTER SET latin1 NOT NULL,
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data bancho_private_messages: ~0 rows (approximately)
/*!40000 ALTER TABLE `bancho_private_messages` DISABLE KEYS */;
/*!40000 ALTER TABLE `bancho_private_messages` ENABLE KEYS */;

-- Dumping structure bancho_settings
CREATE TABLE IF NOT EXISTS `bancho_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `value_int` int(11) NOT NULL DEFAULT '0',
  `value_string` varchar(512) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- Dumping data bancho_settings: ~8 rows (approximately)
/*!40000 ALTER TABLE `bancho_settings` DISABLE KEYS */;
INSERT INTO `bancho_settings` (`id`, `name`, `value_int`, `value_string`) VALUES
	(1, 'bancho_maintenance', 0, ''),
	(2, 'free_direct', 1, ''),
	(3, 'menu_icon', 1, 'https://i.ppy.sh/logo.png | https://ripple.moe'),
	(4, 'login_messages', 1, ''),
	(5, 'restricted_joke', 0, 'You\'re banned from the server.'),
	(6, 'login_notification', 1, 'Welcome to Light Ripple!'),
	(7, 'osu_versions', 0, ''),
	(8, 'osu_md5s', 0, '');
/*!40000 ALTER TABLE `bancho_settings` ENABLE KEYS */;

-- Dumping structure bancho_tokens
CREATE TABLE IF NOT EXISTS `bancho_tokens` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `token` varchar(16) CHARACTER SET latin1 NOT NULL,
  `osu_id` int(11) NOT NULL,
  `latest_message_id` int(11) NOT NULL,
  `latest_private_message_id` int(11) NOT NULL,
  `latest_packet_time` int(11) NOT NULL,
  `latest_heavy_packet_time` int(11) NOT NULL,
  `joined_channels` varchar(512) CHARACTER SET latin1 NOT NULL,
  `game_mode` tinyint(4) NOT NULL,
  `action` int(11) NOT NULL,
  `action_text` varchar(128) CHARACTER SET latin1 NOT NULL,
  `kicked` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data bancho_tokens: ~0 rows (approximately)
/*!40000 ALTER TABLE `bancho_tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `bancho_tokens` ENABLE KEYS */;

-- Dumping structure beatmaps
CREATE TABLE IF NOT EXISTS `beatmaps` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `beatmap_id` int(11) NOT NULL DEFAULT '0',
  `beatmapset_id` int(11) NOT NULL DEFAULT '0',
  `beatmap_md5` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `song_name` text CHARACTER SET latin1 NOT NULL,
  `ar` float NOT NULL DEFAULT '0',
  `od` float NOT NULL DEFAULT '0',
  `mode` int(11) NOT NULL DEFAULT '0',
  `rating` int(11) NOT NULL DEFAULT '10',
  `difficulty_std` float NOT NULL DEFAULT '0',
  `difficulty_taiko` float NOT NULL DEFAULT '0',
  `difficulty_ctb` float NOT NULL DEFAULT '0',
  `difficulty_mania` float NOT NULL DEFAULT '0',
  `max_combo` int(11) NOT NULL DEFAULT '0',
  `hit_length` int(11) NOT NULL DEFAULT '0',
  `bpm` int(11) NOT NULL DEFAULT '0',
  `playcount` int(11) NOT NULL DEFAULT '0',
  `passcount` int(11) NOT NULL DEFAULT '0',
  `ranked` tinyint(4) NOT NULL DEFAULT '0',
  `latest_update` int(11) NOT NULL DEFAULT '0',
  `ranked_status_freezed` tinyint(1) NOT NULL DEFAULT '0',
  `pp_100` int(11) NOT NULL DEFAULT '0',
  `pp_99` int(11) NOT NULL DEFAULT '0',
  `pp_98` int(11) NOT NULL DEFAULT '0',
  `pp_95` int(11) NOT NULL DEFAULT '0',
  `disable_pp` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `id` (`id`),
  KEY `id_2` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data beatmaps: ~0 rows (approximately)
/*!40000 ALTER TABLE `beatmaps` DISABLE KEYS */;
/*!40000 ALTER TABLE `beatmaps` ENABLE KEYS */;

-- Dumping structure beatmaps_names
CREATE TABLE IF NOT EXISTS `beatmaps_names` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `beatmap_md5` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `beatmap_name` varchar(256) CHARACTER SET latin1 NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data beatmaps_names: ~0 rows (approximately)
/*!40000 ALTER TABLE `beatmaps_names` DISABLE KEYS */;
/*!40000 ALTER TABLE `beatmaps_names` ENABLE KEYS */;

-- Dumping structure beatmaps_rating
CREATE TABLE IF NOT EXISTS `beatmaps_rating` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `beatmap_md5` varchar(32) NOT NULL,
  `rating` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data beatmaps_rating: ~0 rows (approximately)
/*!40000 ALTER TABLE `beatmaps_rating` DISABLE KEYS */;
/*!40000 ALTER TABLE `beatmaps_rating` ENABLE KEYS */;

-- Dumping structure clans
CREATE TABLE IF NOT EXISTS `clans` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `description` text NOT NULL,
  `icon` text NOT NULL,
  `tag` varchar(6) NOT NULL,
  `mlimit` int(11) NOT NULL DEFAULT '16',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data clans: ~0 rows (approximately)
/*!40000 ALTER TABLE `clans` DISABLE KEYS */;
/*!40000 ALTER TABLE `clans` ENABLE KEYS */;

-- Dumping structure clans_invites
CREATE TABLE IF NOT EXISTS `clans_invites` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `clan` int(11) NOT NULL,
  `invite` varchar(8) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data clans_invites: ~0 rows (approximately)
/*!40000 ALTER TABLE `clans_invites` DISABLE KEYS */;
/*!40000 ALTER TABLE `clans_invites` ENABLE KEYS */;

-- Dumping structure comments
CREATE TABLE IF NOT EXISTS `comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `beatmap_id` int(11) NOT NULL DEFAULT '0',
  `beatmapset_id` int(11) NOT NULL DEFAULT '0',
  `score_id` int(11) NOT NULL DEFAULT '0',
  `mode` int(11) NOT NULL,
  `comment` varchar(128) CHARACTER SET utf8 NOT NULL,
  `time` int(11) NOT NULL,
  `who` varchar(11) NOT NULL,
  `special_format` varchar(2556) CHARACTER SET utf8 DEFAULT 'FFFFFF',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Dumping data comments: ~0 rows (approximately)
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;

-- Dumping structure discord_roles
CREATE TABLE IF NOT EXISTS `discord_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  `discordid` bigint(19) NOT NULL,
  `roleid` bigint(19) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data discord_roles: ~0 rows (approximately)
/*!40000 ALTER TABLE `discord_roles` DISABLE KEYS */;
/*!40000 ALTER TABLE `discord_roles` ENABLE KEYS */;

-- Dumping structure docs
CREATE TABLE IF NOT EXISTS `docs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `doc_name` varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT 'New Documentation File',
  `doc_contents` mediumtext CHARACTER SET latin1 NOT NULL,
  `public` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `old_name` varchar(200) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- Dumping data docs: 0 rows
/*!40000 ALTER TABLE `docs` DISABLE KEYS */;
/*!40000 ALTER TABLE `docs` ENABLE KEYS */;

-- Dumping structure hw_user
CREATE TABLE IF NOT EXISTS `hw_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  `mac` varchar(32) CHARACTER SET latin1 NOT NULL,
  `unique_id` varchar(32) CHARACTER SET latin1 NOT NULL,
  `disk_id` varchar(32) CHARACTER SET latin1 NOT NULL,
  `occurencies` int(11) NOT NULL DEFAULT '0',
  `activated` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data hw_user: ~0 rows (approximately)
/*!40000 ALTER TABLE `hw_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `hw_user` ENABLE KEYS */;

-- Dumping structure identity_tokens
CREATE TABLE IF NOT EXISTS `identity_tokens` (
  `userid` int(11) NOT NULL,
  `token` varchar(64) CHARACTER SET latin1 NOT NULL,
  UNIQUE KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data identity_tokens: ~0 rows (approximately)
/*!40000 ALTER TABLE `identity_tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `identity_tokens` ENABLE KEYS */;

-- Dumping structure ip_user
CREATE TABLE IF NOT EXISTS `ip_user` (
  `userid` int(11) NOT NULL,
  `ip` text CHARACTER SET latin1 NOT NULL,
  `occurencies` int(11) NOT NULL,
  PRIMARY KEY (`userid`),
  UNIQUE KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data ip_user: ~0 rows (approximately)
/*!40000 ALTER TABLE `ip_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `ip_user` ENABLE KEYS */;

-- Dumping structure irc_tokens
CREATE TABLE IF NOT EXISTS `irc_tokens` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `token` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  UNIQUE KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data irc_tokens: ~0 rows (approximately)
/*!40000 ALTER TABLE `irc_tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `irc_tokens` ENABLE KEYS */;

-- Dumping structure leaderboard_ctb
CREATE TABLE IF NOT EXISTS `leaderboard_ctb` (
  `position` int(10) unsigned NOT NULL,
  `user` int(11) NOT NULL,
  `v` bigint(20) NOT NULL,
  PRIMARY KEY (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data leaderboard_ctb: ~0 rows (approximately)
/*!40000 ALTER TABLE `leaderboard_ctb` DISABLE KEYS */;
/*!40000 ALTER TABLE `leaderboard_ctb` ENABLE KEYS */;

-- Dumping structure leaderboard_mania
CREATE TABLE IF NOT EXISTS `leaderboard_mania` (
  `position` int(10) unsigned NOT NULL,
  `user` int(11) NOT NULL,
  `v` bigint(20) NOT NULL,
  PRIMARY KEY (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data leaderboard_mania: ~0 rows (approximately)
/*!40000 ALTER TABLE `leaderboard_mania` DISABLE KEYS */;
/*!40000 ALTER TABLE `leaderboard_mania` ENABLE KEYS */;

-- Dumping structure leaderboard_std
CREATE TABLE IF NOT EXISTS `leaderboard_std` (
  `position` int(10) unsigned NOT NULL,
  `user` int(11) NOT NULL,
  `v` bigint(20) NOT NULL,
  PRIMARY KEY (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data leaderboard_std: ~0 rows (approximately)
/*!40000 ALTER TABLE `leaderboard_std` DISABLE KEYS */;
/*!40000 ALTER TABLE `leaderboard_std` ENABLE KEYS */;

-- Dumping structure leaderboard_taiko
CREATE TABLE IF NOT EXISTS `leaderboard_taiko` (
  `position` int(10) unsigned NOT NULL,
  `user` int(11) NOT NULL,
  `v` bigint(20) NOT NULL,
  PRIMARY KEY (`position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data leaderboard_taiko: ~0 rows (approximately)
/*!40000 ALTER TABLE `leaderboard_taiko` DISABLE KEYS */;
/*!40000 ALTER TABLE `leaderboard_taiko` ENABLE KEYS */;

-- Dumping structure main_menu_icons
CREATE TABLE IF NOT EXISTS `main_menu_icons` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `is_current` int(11) NOT NULL,
  `file_id` varchar(128) NOT NULL,
  `name` varchar(256) NOT NULL,
  `url` text CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

-- Dumping data main_menu_icons: ~1 rows (approximately)
/*!40000 ALTER TABLE `main_menu_icons` DISABLE KEYS */;
INSERT INTO `main_menu_icons` (`id`, `is_current`, `file_id`, `name`, `url`) VALUES
	(1, 1, 'logo', 'Ainu!', 'https://ainu.pw/');
/*!40000 ALTER TABLE `main_menu_icons` ENABLE KEYS */;

-- Dumping structure osin_access
CREATE TABLE IF NOT EXISTS `osin_access` (
  `scope` int(11) NOT NULL DEFAULT '0',
  `created_at` int(11) NOT NULL DEFAULT '0',
  `client` int(11) NOT NULL DEFAULT '0',
  `extra` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data osin_access: ~0 rows (approximately)
/*!40000 ALTER TABLE `osin_access` DISABLE KEYS */;
/*!40000 ALTER TABLE `osin_access` ENABLE KEYS */;

-- Dumping structure osin_client
CREATE TABLE IF NOT EXISTS `osin_client` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `secret` varchar(64) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `extra` varchar(127) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `redirect_uri` varchar(127) CHARACTER SET latin1 NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data osin_client: ~0 rows (approximately)
/*!40000 ALTER TABLE `osin_client` DISABLE KEYS */;
/*!40000 ALTER TABLE `osin_client` ENABLE KEYS */;

-- Dumping structure osin_client_user
CREATE TABLE IF NOT EXISTS `osin_client_user` (
  `client_id` int(11) NOT NULL DEFAULT '0',
  `user` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data osin_client_user: ~0 rows (approximately)
/*!40000 ALTER TABLE `osin_client_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `osin_client_user` ENABLE KEYS */;

-- Dumping structure password_recovery
CREATE TABLE IF NOT EXISTS `password_recovery` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `k` varchar(80) CHARACTER SET latin1 NOT NULL,
  `u` varchar(30) CHARACTER SET latin1 NOT NULL,
  `t` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- Dumping data password_recovery: 0 rows
/*!40000 ALTER TABLE `password_recovery` DISABLE KEYS */;
/*!40000 ALTER TABLE `password_recovery` ENABLE KEYS */;

-- Dumping structure privileges_groups
CREATE TABLE IF NOT EXISTS `privileges_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET latin1 NOT NULL,
  `privileges` int(11) NOT NULL,
  `color` varchar(32) CHARACTER SET latin1 NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;

-- Dumping data privileges_groups: ~13 rows (approximately)
/*!40000 ALTER TABLE `privileges_groups` DISABLE KEYS */;
INSERT INTO `privileges_groups` (`id`, `name`, `privileges`, `color`) VALUES
	(1, 'Banned', 0, ''),
	(2, 'BAT', 267, ''),
	(3, 'Chat Moderators', 2883911, 'success'),
	(4, 'Admin', 1048575, 'danger'),
	(5, 'Developer', 1043995, 'info'),
	(6, 'Donor', 7, 'warning'),
	(7, 'God', 1048575, 'info'),
	(8, 'Normal User', 3, 'primary'),
	(9, 'Pending', 1048576, 'default'),
	(10, 'Restricted', 2, ''),
	(11, 'Beatmap Nominator', 267, 'primary'),
	(12, 'Owner', 3145727, 'info'),
	(13, 'Community Manager', 918015, 'success');
/*!40000 ALTER TABLE `privileges_groups` ENABLE KEYS */;

-- Dumping structure profile_backgrounds
CREATE TABLE IF NOT EXISTS `profile_backgrounds` (
  `uid` int(11) NOT NULL,
  `time` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `value` text CHARACTER SET latin1 NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data profile_backgrounds: ~0 rows (approximately)
/*!40000 ALTER TABLE `profile_backgrounds` DISABLE KEYS */;
/*!40000 ALTER TABLE `profile_backgrounds` ENABLE KEYS */;

-- Dumping structure rank_requests
CREATE TABLE IF NOT EXISTS `rank_requests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  `bid` int(11) NOT NULL,
  `type` varchar(8) CHARACTER SET latin1 NOT NULL,
  `time` int(11) NOT NULL,
  `blacklisted` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bid` (`bid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data rank_requests: ~0 rows (approximately)
/*!40000 ALTER TABLE `rank_requests` DISABLE KEYS */;
/*!40000 ALTER TABLE `rank_requests` ENABLE KEYS */;

-- Dumping structure rap_logs
CREATE TABLE IF NOT EXISTS `rap_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  `text` text NOT NULL,
  `datetime` int(11) NOT NULL,
  `through` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data rap_logs: ~0 rows (approximately)
/*!40000 ALTER TABLE `rap_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `rap_logs` ENABLE KEYS */;

-- Dumping structure remember
CREATE TABLE IF NOT EXISTS `remember` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` int(11) NOT NULL,
  `series_identifier` int(11) DEFAULT NULL,
  `token_sha` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data remember: ~0 rows (approximately)
/*!40000 ALTER TABLE `remember` DISABLE KEYS */;
/*!40000 ALTER TABLE `remember` ENABLE KEYS */;

-- Dumping structure reports
CREATE TABLE IF NOT EXISTS `reports` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `from_uid` int(11) NOT NULL,
  `to_uid` int(11) NOT NULL,
  `reason` text NOT NULL,
  `chatlog` text NOT NULL,
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data reports: ~0 rows (approximately)
/*!40000 ALTER TABLE `reports` DISABLE KEYS */;
/*!40000 ALTER TABLE `reports` ENABLE KEYS */;

-- Dumping structure rx_beatmap_playcount
CREATE TABLE IF NOT EXISTS `rx_beatmap_playcount` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `beatmap_id` int(11) DEFAULT NULL,
  `game_mode` int(11) DEFAULT NULL,
  `playcount` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `playcount_index` (`user_id`,`beatmap_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data rx_beatmap_playcount: ~0 rows (approximately)
/*!40000 ALTER TABLE `rx_beatmap_playcount` DISABLE KEYS */;
/*!40000 ALTER TABLE `rx_beatmap_playcount` ENABLE KEYS */;

-- Dumping structure rx_stats
CREATE TABLE IF NOT EXISTS `rx_stats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) CHARACTER SET latin1 NOT NULL,
  `username_aka` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `user_color` varchar(16) CHARACTER SET latin1 NOT NULL DEFAULT 'black',
  `user_style` varchar(128) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `favourite_mode` int(11) NOT NULL DEFAULT '0',
  `level_std` int(11) NOT NULL DEFAULT '1',
  `level_taiko` int(11) NOT NULL DEFAULT '1',
  `level_mania` int(11) NOT NULL DEFAULT '1',
  `level_ctb` int(11) NOT NULL DEFAULT '1',
  `total_score_std` int(11) NOT NULL DEFAULT '0',
  `total_score_taiko` int(11) NOT NULL DEFAULT '0',
  `total_score_mania` int(11) NOT NULL DEFAULT '0',
  `total_score_ctb` int(11) NOT NULL DEFAULT '0',
  `total_hits_std` int(11) NOT NULL DEFAULT '0',
  `total_hits_taiko` int(11) NOT NULL DEFAULT '0',
  `total_hits_ctb` int(11) NOT NULL DEFAULT '0',
  `total_hits_mania` int(11) NOT NULL DEFAULT '0',
  `playtime_std` int(11) NOT NULL DEFAULT '0',
  `playtime_taiko` int(11) NOT NULL DEFAULT '0',
  `playtime_mania` int(11) NOT NULL DEFAULT '0',
  `playtime_ctb` int(11) NOT NULL DEFAULT '0',
  `ranked_score_std` bigint(11) NOT NULL DEFAULT '0',
  `ranked_score_taiko` int(11) NOT NULL DEFAULT '0',
  `ranked_score_mania` int(11) NOT NULL DEFAULT '0',
  `ranked_score_ctb` int(11) NOT NULL DEFAULT '0',
  `avg_accuracy_std` double NOT NULL DEFAULT '0',
  `avg_accuracy_taiko` double NOT NULL DEFAULT '0',
  `avg_accuracy_mania` double NOT NULL DEFAULT '0',
  `avg_accuracy_ctb` double NOT NULL DEFAULT '0',
  `playcount_std` int(11) NOT NULL DEFAULT '0',
  `playcount_taiko` int(11) NOT NULL DEFAULT '0',
  `playcount_mania` int(11) NOT NULL DEFAULT '0',
  `playcount_ctb` int(11) NOT NULL DEFAULT '0',
  `pp_std` int(11) NOT NULL DEFAULT '0',
  `pp_mania` int(11) NOT NULL DEFAULT '0',
  `pp_ctb` int(11) NOT NULL DEFAULT '0',
  `pp_taiko` int(11) NOT NULL DEFAULT '0',
  `country` char(2) NOT NULL DEFAULT 'XX',
  `unrestricted_pp` int(11) NOT NULL DEFAULT '0',
  `ppboard` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8;

-- Dumping data rx_stats: ~0 rows (approximately)
/*!40000 ALTER TABLE `rx_stats` DISABLE KEYS */;
INSERT INTO `rx_stats` (`id`, `username`, `username_aka`, `user_color`, `user_style`, `favourite_mode`, `level_std`, `level_taiko`, `level_mania`, `level_ctb`, `total_score_std`, `total_score_taiko`, `total_score_mania`, `total_score_ctb`, `total_hits_std`, `total_hits_taiko`, `total_hits_ctb`, `total_hits_mania`, `playtime_std`, `playtime_taiko`, `playtime_mania`, `playtime_ctb`, `ranked_score_std`, `ranked_score_taiko`, `ranked_score_mania`, `ranked_score_ctb`, `avg_accuracy_std`, `avg_accuracy_taiko`, `avg_accuracy_mania`, `avg_accuracy_ctb`, `playcount_std`, `playcount_taiko`, `playcount_mania`, `playcount_ctb`, `pp_std`, `pp_mania`, `pp_ctb`, `pp_taiko`, `country`, `unrestricted_pp`, `ppboard`) VALUES
	(999, 'AC', '', 'black', '', 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'JP', 1, 1);
/*!40000 ALTER TABLE `rx_stats` ENABLE KEYS */;

-- Dumping structure scores
CREATE TABLE IF NOT EXISTS `scores` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `beatmap_md5` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `userid` int(11) NOT NULL,
  `score` bigint(20) DEFAULT NULL,
  `max_combo` int(11) NOT NULL DEFAULT '0',
  `full_combo` tinyint(1) NOT NULL DEFAULT '0',
  `mods` int(11) NOT NULL DEFAULT '0',
  `300_count` int(11) NOT NULL DEFAULT '0',
  `100_count` int(11) NOT NULL DEFAULT '0',
  `50_count` int(11) NOT NULL DEFAULT '0',
  `katus_count` int(11) NOT NULL DEFAULT '0',
  `gekis_count` int(11) NOT NULL DEFAULT '0',
  `misses_count` int(11) NOT NULL DEFAULT '0',
  `time` varchar(18) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `play_mode` tinyint(4) NOT NULL DEFAULT '0',
  `completed` tinyint(11) NOT NULL DEFAULT '0',
  `accuracy` float(15,12) DEFAULT NULL,
  `pp` double DEFAULT '0',
  `playtime` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data scores: ~0 rows (approximately)
/*!40000 ALTER TABLE `scores` DISABLE KEYS */;
/*!40000 ALTER TABLE `scores` ENABLE KEYS */;

-- Dumping structure scores_relax
CREATE TABLE IF NOT EXISTS `scores_relax` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `beatmap_md5` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `userid` int(11) NOT NULL,
  `score` bigint(20) DEFAULT NULL,
  `max_combo` int(11) NOT NULL DEFAULT '0',
  `full_combo` tinyint(1) NOT NULL DEFAULT '0',
  `mods` int(11) NOT NULL DEFAULT '0',
  `300_count` int(11) NOT NULL DEFAULT '0',
  `100_count` int(11) NOT NULL DEFAULT '0',
  `50_count` int(11) NOT NULL DEFAULT '0',
  `katus_count` int(11) NOT NULL DEFAULT '0',
  `gekis_count` int(11) NOT NULL DEFAULT '0',
  `misses_count` int(11) NOT NULL DEFAULT '0',
  `time` varchar(18) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `play_mode` tinyint(4) NOT NULL DEFAULT '0',
  `completed` tinyint(11) NOT NULL DEFAULT '0',
  `accuracy` float(15,12) DEFAULT NULL,
  `pp` double DEFAULT '0',
  `playtime` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1073741824 DEFAULT CHARSET=utf8;

-- Dumping data scores_relax: ~0 rows (approximately)
/*!40000 ALTER TABLE `scores_relax` DISABLE KEYS */;
/*!40000 ALTER TABLE `scores_relax` ENABLE KEYS */;

-- Dumping structure system_settings
CREATE TABLE IF NOT EXISTS `system_settings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET latin1 NOT NULL,
  `value_int` int(11) NOT NULL DEFAULT '0',
  `value_string` varchar(512) CHARACTER SET latin1 NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- Dumping data system_settings: ~10 rows (approximately)
/*!40000 ALTER TABLE `system_settings` DISABLE KEYS */;
INSERT INTO `system_settings` (`id`, `name`, `value_int`, `value_string`) VALUES
	(1, 'website_maintenance', 0, ''),
	(2, 'game_maintenance', 0, ''),
	(3, 'website_global_alert', 0, ''),
	(4, 'website_home_alert', 0, ''),
	(5, 'registrations_enabled', 1, ''),
	(6, 'ccreation_enabled', 1, '');
/*!40000 ALTER TABLE `system_settings` ENABLE KEYS */;

-- Dumping structure tokens
CREATE TABLE IF NOT EXISTS `tokens` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(31) CHARACTER SET latin1 NOT NULL,
  `privileges` int(11) NOT NULL,
  `description` varchar(255) CHARACTER SET latin1 NOT NULL,
  `token` varchar(127) CHARACTER SET latin1 NOT NULL,
  `private` tinyint(4) NOT NULL,
  `last_updated` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data tokens: ~0 rows (approximately)
/*!40000 ALTER TABLE `tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `tokens` ENABLE KEYS */;

-- Dumping structure users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(15) NOT NULL AUTO_INCREMENT,
  `osuver` varchar(256) DEFAULT NULL,
  `username` varchar(30) CHARACTER SET latin1 NOT NULL,
  `username_safe` varchar(30) CHARACTER SET latin1 NOT NULL,
  `ban_datetime` varchar(30) NOT NULL DEFAULT '0',
  `password_md5` varchar(127) CHARACTER SET latin1 NOT NULL,
  `salt` varchar(32) CHARACTER SET latin1 NOT NULL,
  `email` varchar(254) CHARACTER SET latin1 NOT NULL,
  `register_datetime` int(10) NOT NULL,
  `rank` tinyint(1) NOT NULL DEFAULT '1',
  `allowed` tinyint(1) NOT NULL DEFAULT '1',
  `latest_activity` int(10) NOT NULL DEFAULT '0',
  `silence_end` int(11) NOT NULL DEFAULT '0',
  `silence_reason` varchar(127) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `password_version` tinyint(4) NOT NULL DEFAULT '1',
  `privileges` bigint(11) NOT NULL,
  `donor_expire` int(11) NOT NULL DEFAULT '0',
  `flags` int(11) NOT NULL DEFAULT '0',
  `achievements_version` int(11) NOT NULL DEFAULT '4',
  `achievements_0` int(11) NOT NULL DEFAULT '1',
  `achievements_1` int(11) NOT NULL DEFAULT '1',
  `notes` text CHARACTER SET latin1,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8;

-- Dumping data users: ~1 rows (approximately)
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `osuver`, `username`, `username_safe`, `ban_datetime`, `password_md5`, `salt`, `email`, `register_datetime`, `rank`, `allowed`, `latest_activity`, `silence_end`, `silence_reason`, `password_version`, `privileges`, `donor_expire`, `flags`, `achievements_version`, `achievements_0`, `achievements_1`, `notes`) VALUES
	(999, NULL, 'AC', 'ac', '0', '$2y$10$LWyP3VsgCnr7BxEjq33LqO5Wx.ZQUfFR4ZGCzsiShB2FZ3UcXpOXW', '', 'fo@kab.ot', 1566228790, 4, 1, 1569775752, 0, '', 1, 3145727, 2147483647, 0, 0, 1, 1, '');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;

-- Dumping structure users_achievements
CREATE TABLE IF NOT EXISTS `users_achievements` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `achievement_id` int(11) NOT NULL,
  `time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Dumping data users_achievements: ~0 rows (approximately)
/*!40000 ALTER TABLE `users_achievements` DISABLE KEYS */;
/*!40000 ALTER TABLE `users_achievements` ENABLE KEYS */;

-- Dumping structure users_beatmap_playcount
CREATE TABLE IF NOT EXISTS `users_beatmap_playcount` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `beatmap_id` int(11) DEFAULT NULL,
  `game_mode` int(11) DEFAULT NULL,
  `playcount` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `playcount_index` (`user_id`,`beatmap_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data users_beatmap_playcount: ~0 rows (approximately)
/*!40000 ALTER TABLE `users_beatmap_playcount` DISABLE KEYS */;
/*!40000 ALTER TABLE `users_beatmap_playcount` ENABLE KEYS */;

-- Dumping structure users_relationships
CREATE TABLE IF NOT EXISTS `users_relationships` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user1` int(11) NOT NULL,
  `user2` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data users_relationships: ~0 rows (approximately)
/*!40000 ALTER TABLE `users_relationships` DISABLE KEYS */;
/*!40000 ALTER TABLE `users_relationships` ENABLE KEYS */;

-- Dumping structure users_stats
CREATE TABLE IF NOT EXISTS `users_stats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) CHARACTER SET latin1 NOT NULL,
  `username_aka` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `user_color` varchar(16) CHARACTER SET latin1 NOT NULL DEFAULT 'black',
  `user_style` varchar(128) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `ranked_score_std` bigint(20) DEFAULT '0',
  `playcount_std` int(11) NOT NULL DEFAULT '0',
  `total_score_std` bigint(20) DEFAULT '0',
  `replays_watched_std` int(11) unsigned NOT NULL DEFAULT '0',
  `ranked_score_taiko` bigint(20) DEFAULT '0',
  `playcount_taiko` int(11) NOT NULL DEFAULT '0',
  `total_score_taiko` bigint(20) DEFAULT '0',
  `replays_watched_taiko` int(11) NOT NULL DEFAULT '0',
  `ranked_score_ctb` bigint(20) DEFAULT '0',
  `playcount_ctb` int(11) NOT NULL DEFAULT '0',
  `total_score_ctb` bigint(20) DEFAULT '0',
  `replays_watched_ctb` int(11) NOT NULL DEFAULT '0',
  `ranked_score_mania` bigint(20) DEFAULT '0',
  `playcount_mania` int(11) NOT NULL DEFAULT '0',
  `total_score_mania` bigint(20) DEFAULT '0',
  `replays_watched_mania` int(10) unsigned NOT NULL DEFAULT '0',
  `total_hits_std` int(11) NOT NULL DEFAULT '0',
  `total_hits_taiko` int(11) NOT NULL DEFAULT '0',
  `total_hits_ctb` int(11) NOT NULL DEFAULT '0',
  `total_hits_mania` int(11) NOT NULL DEFAULT '0',
  `country` char(2) CHARACTER SET latin1 NOT NULL DEFAULT 'XX',
  `unrestricted_pp` int(11) NOT NULL DEFAULT '0',
  `ppboard` int(11) NOT NULL DEFAULT '0',
  `show_country` tinyint(4) NOT NULL DEFAULT '1',
  `level_std` int(11) NOT NULL DEFAULT '1',
  `level_taiko` int(11) NOT NULL DEFAULT '1',
  `level_ctb` int(11) NOT NULL DEFAULT '1',
  `level_mania` int(11) NOT NULL DEFAULT '1',
  `playtime_std` int(11) NOT NULL DEFAULT '0',
  `playtime_taiko` int(11) NOT NULL DEFAULT '0',
  `playtime_ctb` int(11) NOT NULL DEFAULT '0',
  `playtime_mania` int(11) NOT NULL DEFAULT '0',
  `avg_accuracy_std` float(15,12) NOT NULL DEFAULT '0.000000000000',
  `avg_accuracy_taiko` float(15,12) NOT NULL DEFAULT '0.000000000000',
  `avg_accuracy_ctb` float(15,12) NOT NULL DEFAULT '0.000000000000',
  `avg_accuracy_mania` float(15,12) NOT NULL DEFAULT '0.000000000000',
  `pp_std` int(11) NOT NULL DEFAULT '0',
  `pp_taiko` int(11) NOT NULL DEFAULT '0',
  `pp_ctb` int(11) NOT NULL DEFAULT '0',
  `pp_mania` int(11) NOT NULL DEFAULT '0',
  `badges_shown` varchar(24) CHARACTER SET latin1 NOT NULL DEFAULT '1,0,0,0,0,0',
  `safe_title` tinyint(4) NOT NULL DEFAULT '0',
  `userpage_content` mediumtext CHARACTER SET latin1,
  `play_style` smallint(6) NOT NULL DEFAULT '0',
  `favourite_mode` tinyint(4) NOT NULL DEFAULT '0',
  `prefer_relax` int(11) NOT NULL DEFAULT '0',
  `custom_badge_icon` varchar(32) CHARACTER SET latin1 NOT NULL DEFAULT '',
  `custom_badge_name` varchar(256) NOT NULL DEFAULT '',
  `can_custom_badge` tinyint(1) NOT NULL DEFAULT '0',
  `show_custom_badge` tinyint(1) NOT NULL DEFAULT '0',
  `current_status` varchar(20000) NOT NULL DEFAULT 'Offline',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8;

-- Dumping data users_stats: ~1 rows (approximately)
/*!40000 ALTER TABLE `users_stats` DISABLE KEYS */;
INSERT INTO `users_stats` (`id`, `username`, `username_aka`, `user_color`, `user_style`, `ranked_score_std`, `playcount_std`, `total_score_std`, `replays_watched_std`, `ranked_score_taiko`, `playcount_taiko`, `total_score_taiko`, `replays_watched_taiko`, `ranked_score_ctb`, `playcount_ctb`, `total_score_ctb`, `replays_watched_ctb`, `ranked_score_mania`, `playcount_mania`, `total_score_mania`, `replays_watched_mania`, `total_hits_std`, `total_hits_taiko`, `total_hits_ctb`, `total_hits_mania`, `country`, `unrestricted_pp`, `ppboard`, `show_country`, `level_std`, `level_taiko`, `level_ctb`, `level_mania`, `playtime_std`, `playtime_taiko`, `playtime_ctb`, `playtime_mania`, `avg_accuracy_std`, `avg_accuracy_taiko`, `avg_accuracy_ctb`, `avg_accuracy_mania`, `pp_std`, `pp_taiko`, `pp_ctb`, `pp_mania`, `badges_shown`, `safe_title`, `userpage_content`, `play_style`, `favourite_mode`, `prefer_relax`, `custom_badge_icon`, `custom_badge_name`, `can_custom_badge`, `show_custom_badge`, `current_status`) VALUES
	(999, 'AC', '', 'black', '', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'JP', 1, 0, 1, 65535, 1, 1, 1, 0, 0, 0, 0, 100.000000000000, 0.000000000000, 0.000000000000, 0.000000000000, 0, 0, 0, 0, '3,4,11,0,0,0', 0, '', 0, 0, 0, '', '', 1, 1, 'Dead');
/*!40000 ALTER TABLE `users_stats` ENABLE KEYS */;

-- Dumping structure user_badges
CREATE TABLE IF NOT EXISTS `user_badges` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL,
  `badge` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data user_badges: ~0 rows (approximately)
/*!40000 ALTER TABLE `user_badges` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_badges` ENABLE KEYS */;

-- Dumping structure user_clans
CREATE TABLE IF NOT EXISTS `user_clans` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL,
  `clan` int(11) NOT NULL,
  `perms` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dumping data user_clans: ~0 rows (approximately)
/*!40000 ALTER TABLE `user_clans` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_clans` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
