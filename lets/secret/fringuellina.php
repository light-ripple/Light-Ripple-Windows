<?php
//echo 'Its been a while since ive done php. Let me hide my gun first.';

class Fringuellina {
	public static $cakeRecipeName = "Cakes"; 

    public static function PrintPage(){

		if ($_GET['id'] && !empty($_GET['id'])) {
			Fringuellina::PrintCake();
			return;
		}

		if ($_GET['uid'] && !empty($_GET['uid'])) {
			Fringuellina::PrintUserCakes();
			return;
		}

        // Multiple pages
		$pageInterval = 100;
		$from = (isset($_GET["from"])) ? $_GET["from"] : 999;
		$to = $from+$pageInterval;
		$users = $GLOBALS['db']->fetchAll('SELECT * FROM users WHERE id >= ? AND id < ?', [$from, $to]);
		$groups = $GLOBALS["db"]->fetchAll("SELECT * FROM privileges_groups");

        // Print stuff
		echo '<div id="wrapper">';
		printAdminSidebar();
		echo '<div id="page-content-wrapper">';
		// Maintenance check
		P::MaintenanceStuff();
		// Print Success if set
		if (isset($_GET['s']) && !empty($_GET['s'])) {
			P::SuccessMessageStaccah($_GET['s']);
		}
		// Print Exception if set
		if (isset($_GET['e']) && !empty($_GET['e'])) {
			P::ExceptionMessageStaccah($_GET['e']);
        }
        
		// Get values
		//$wm = current($GLOBALS['db']->fetch("SELECT value_int FROM system_settings WHERE name = 'website_maintenance'"));
		
		echo '<p align="center"><font size=5><i class="fa fa-birthday-cake"></i>	Cakes</font></p><br>';
		// Quick edit/silence/kick user button
		echo '<p align="center"><button type="button" class="btn btn-primary" data-toggle="modal" data-target="#quickLookupUserModal">Quick lookup user (username)</button>';
		echo '&nbsp;&nbsp; <button type="button" class="btn btn-warning" data-toggle="modal" data-target="#quickLookupScoreidModal">Quick lookup cake (score_id)</button>';
		echo '</p>';
		// Users plays table
		echo '<table class="table table-striped table-hover table-50-center">
		<thead>
		<tr><th class="text-center"><i class="fa fa-user"></i>	ID</th><th class="text-center">Username</th><th class="text-center">Cakes</th><th class="text-center">Bad cakes</th><th class="text-center">Bad flags</th><th class="text-center">Status</th><th class="text-center">Actions</th></tr>
		</thead>
		<tbody>';
		foreach ($users as $user) {
			// Get allowed color/text
			$statusColor = "success";
			$statusText = "Ok";
			if (($user["privileges"] & Privileges::UserPublic) == 0 && ($user["privileges"] & Privileges::UserNormal) == 0) {
				// Not visible and not active, banned
				$statusColor = "danger";
				$statusText = "Banned";
			} else if (($user["privileges"] & Privileges::UserPublic) == 0 && ($user["privileges"] & Privileges::UserNormal) > 0) {
				// Not visible but active, restricted
				$statusColor = "warning";
				$statusText = "Restricted";
			} else if (($user["privileges"] & Privileges::UserPublic) > 0 && ($user["privileges"] & Privileges::UserNormal) == 0) {
				// Visible but not active, disabled (not supported yet)
				$statusColor = "default";
				$statusText = "Locked";
            }
            
            $cakes = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ?', [$user["id"]]));

			$badCakes = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ? AND detected NOT LIKE ?', [$user["id"], '[]']));
			$badFlags = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ? AND flags NOT IN (0,4)', [$user["id"]]));

			// Print row
			echo '<tr>';
			echo '<td><p class="text-center">'.$user['id'].'</p></td>';
			echo '<td><p class="text-center"><b>'.$user['username'].'</b></p></td>';
			echo '<td><p class="text-center">'.$cakes.'</p></td>';
			echo '<td><p class="text-center">'.$badCakes.'</p></td>';
			echo '<td><p class="text-center">'.$badFlags.'</p></td>';
            echo '<td><p class="text-center"><span class="label label-'.$statusColor.'">'.$statusText.'</span></p></td>';
            echo '<td><p class="text-center"><div class="btn-group">';
			echo '<a title="Edit user" class="btn btn-xs btn-primary" href="index.php?p=128&uid='.$user['id'].'"><span class="glyphicon glyphicon-pencil"></span></a>';
			if (hasPrivilege(Privileges::AdminBanUsers)) {
				if (isBanned($user["id"])) {
					echo '<a title="Unban user" class="btn btn-xs btn-success" onclick="sure(\'submit.php?action=banUnbanUser&id='.$user['id'].'\')"><span class="glyphicon glyphicon-thumbs-up"></span></a>';
				} else {
					echo '<a title="Ban user" class="btn btn-xs btn-warning" onclick="sure(\'submit.php?action=banUnbanUser&id='.$user['id'].'\')"><span class="glyphicon glyphicon-thumbs-down"></span></a>';
				}
				if (isRestricted($user["id"])) {
					echo '<a title="Remove restrictions" class="btn btn-xs btn-success" onclick="sure(\'submit.php?action=restrictUnrestrictUser&id='.$user['id'].'\')"><span class="glyphicon glyphicon-ok-circle"></span></a>';
				} else {
					echo '<a title="Restrict user" class="btn btn-xs btn-warning" onclick="sure(\'submit.php?action=restrictUnrestrictUser&id='.$user['id'].'\')"><span class="glyphicon glyphicon-remove-circle"></span></a>';
				}
			}
			echo '</div>';
			echo '</td>';
			echo '</tr>';
		}
		echo '</tbody></table>';
		echo '<p align="center"><a href="index.php?p=128&from='.($from-($pageInterval+1)).'">< Previous page</a> | <a href="index.php?p=128&from='.($to).'">Next page ></a></p>';
		echo '</div>';



		Fringuellina::PrintLookupUserModule();
		Fringuellina::PrintLookupCakeModule();
	}
	
	public static function PrintUserCakes(){
		// Print stuff
		echo '<div id="wrapper">';
		printAdminSidebar();
		echo '<div id="page-content-wrapper">';
		// Maintenance check
		P::MaintenanceStuff();
		// Print Success if set
		if (isset($_GET['s']) && !empty($_GET['s'])) {
			P::SuccessMessageStaccah($_GET['s']);
		}
		// Print Exception if set
		if (isset($_GET['e']) && !empty($_GET['e'])) {
			P::ExceptionMessageStaccah($_GET['e']);
		}

		$uid = $_GET['uid'];

		$user = $GLOBALS['db']->fetch('SELECT * FROM users WHERE id = ?', [$uid]);

		$statusColor = "success";
		$statusText = "Ok";
		if (($user["privileges"] & Privileges::UserPublic) == 0 && ($user["privileges"] & Privileges::UserNormal) == 0) {
			// Not visible and not active, banned
			$statusColor = "danger";
			$statusText = "Banned";
		} else if (($user["privileges"] & Privileges::UserPublic) == 0 && ($user["privileges"] & Privileges::UserNormal) > 0) {
			// Not visible but active, restricted
			$statusColor = "warning";
			$statusText = "Restricted";
		} else if (($user["privileges"] & Privileges::UserPublic) > 0 && ($user["privileges"] & Privileges::UserNormal) == 0) {
			// Visible but not active, disabled (not supported yet)
			$statusColor = "default";
			$statusText = "Locked";
		}

		$typeOlSelect = 0;

		if (isset($_GET['q']) && !empty($_GET['q'])) {
			$typeOlSelect = intval($_GET['q']);
		}

		$query = "SELECT * FROM cakes WHERE userid = ?";
		if ($typeOlSelect == 1)
			$query .= " AND detected NOT LIKE '[]'";
		else if ($typeOlSelect == 2)
			$query .= " AND flags NOT IN (0,4)";
		else if ($typeOlSelect == 3)
			$query .= " AND (detected NOT LIKE '[]' OR flags NOT IN (0,4))";

		$order = " ORDER BY id ASC";
		if (isset($_GET['d'])) {
			$order = " ORDER BY id DESC";
		}

		$query .= $order;

		$page = 0;
		if (isset($_GET['l']) && !empty($_GET['l'])){
			$page = intval($_GET['l']);
		}

		$query .= " LIMIT ".$page.", 100";

		$cakes = $GLOBALS['db']->fetchAll($query, [$uid]);

		$cakeCount = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ?', [$uid]));
		$badCakes = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ? AND detected NOT LIKE ?', [$uid, '[]']));
		$badFlags = current($GLOBALS['db']->fetch('SELECT COUNT(*) FROM cakes WHERE userid = ? AND flags NOT IN (0,4)', [$uid]));

		echo '<div class="row">';
		$hrefpage = "index.php?p=128";
		foreach ($_GET as $key => $value){
			if ($key != "q")
				$hrefpage .= "&".$key."=".$value;
		}

		$btn = [1, 2];
		$btncss = ["box-shadow: 0 0 40px #0000ff !important;","",""];

		if ($typeOlSelect == 1){
			$btn = [0, 3];
			$btncss = ["", "box-shadow: 0 0 40px #ff0000 !important;", ""];
		}
		else if($typeOlSelect == 2){
			$btn = [3, 0];
			$btncss = ["", "", "box-shadow: 0 0 40px #ffff00 !important;"];
		}
		else if ($typeOlSelect == 3){
			$btn = [2, 1];
			$btncss = ["", "box-shadow: 0 0 40px #ff0000 !important;", "box-shadow: 0 0 40px #ffff00 !important;"];
		}

		printAdminPanel('primary" onclick=\'window.location.href="'.$hrefpage.'&q=0";\'" style="cursor:pointer;'.$btncss[0], 'fa fa-birthday-cake fa-5x', $cakeCount, 'Cakes');
		printAdminPanel('red" onclick=\'window.location.href="'.$hrefpage.'&q='.$btn[0].'";\'" style="cursor:pointer;'.$btncss[1], 'fa fa-thumbs-down fa-5x', $badCakes, 'Bad cakes');
		printAdminPanel('yellow" onclick=\'window.location.href="'.$hrefpage.'&q='.$btn[1].'";\'" style="cursor:pointer;'.$btncss[2], 'fa fa-flag fa-5x', $badFlags, 'Bad flags');
		printAdminPanel($statusColor, 'fa fa-id-card fa-5x', $statusText, 'Status');
		echo '</div>';

		echo '<p align="center"><font size=5><i class="fa fa-birthday-cake"></i>	'.$user['username'].'\'s Cakes</font></p><br>';

		echo '<table class="table table-striped table-hover table-50-center">
		<thead>
		<tr><th class="text-center"><i class="fa fa-birthday-cake"></i>	Cake ID</th><th class="text-center">Score ID</th><th class="text-center">Cake comment</th><th class="text-center">Flags</th><th class="text-center">Actions</th></tr>
		</thead>
		<tbody>';
		foreach ($cakes as $cake) {
			echo '<td><p class="text-center">'.$cake['id'].'</p></td>';
			echo '<td><p class="text-center">'.$cake['score_id'].'</p></td>';
			echo '<td><p class="text-center">'.$cake['detected'].'</p></td>';
			echo '<td><p class="text-center">'.$cake['flags'].'</p></td>';
			echo '<td><p class="text-center"><a href="index.php?p=129&id='.$cake['id'].'" type="button" class="btn btn-primary">Check cake</a></p></td>';

			echo '</tr>';
		}
		echo '</tbody></table>';

		$hrefpage = "index.php?p=128";
		foreach ($_GET as $key => $value){
			if ($key != "l")
				$hrefpage .= "&".$key."=".$value;
		}

		echo '<p align="center"><a href="'.$hrefpage.'&l='.max($page-100, 0).'">< Previous page</a> | <a href="'.$hrefpage.'&l='.min($page+100, $cakeCount-100).'">Next page ></a></p>';

		echo '</div></div>';

		echo '<script>
		
		</script>';
	}

	public static function PrintInfoPage(){
		//Redirect scoreID to the cakeID
		if ($_GET['sid'] && !empty($_GET['sid'])) {
			$cakeID = $GLOBALS['db']->fetch('SELECT id FROM cakes WHERE score_id = ?', [$_GET['sid']])["id"];
			if ($cakeID != null)
				header('Location: index.php?p=129&id='.$cakeID);
			else
				header('Location: index.php?p=128');
			exit();
		}

		// Print stuff
		echo '<div id="wrapper">';
		printAdminSidebar();
		echo '<div id="page-content-wrapper">';
		// Maintenance check
		P::MaintenanceStuff();
		// Print Success if set
		if (isset($_GET['s']) && !empty($_GET['s'])) {
			P::SuccessMessageStaccah($_GET['s']);
		}
		// Print Exception if set
		if (isset($_GET['e']) && !empty($_GET['e'])) {
			P::ExceptionMessageStaccah($_GET['e']);
		}

		$id = $_GET['id'];

		$cake = $GLOBALS['db']->fetch('SELECT * FROM cakes WHERE id = ?', [$id]);
		$eggs = $GLOBALS['db']->fetchAll('SELECT * FROM eggs');
		$user = $GLOBALS['db']->fetch('SELECT * FROM users WHERE id = ?', [$cake['userid']]);

		$flags = Fringuellina::makeFlagString($cake['flags']);

		$beatmap_md5 = $GLOBALS['db']->fetch('SELECT beatmap_md5 FROM scores WHERE id = ?', [$cake['score_id']])['beatmap_md5'];

		$beatmap = $GLOBALS['db']->fetch('SELECT beatmap_id,beatmapset_id FROM beatmaps WHERE beatmap_md5 = ?', [$beatmap_md5]);

		$cakeCommentsList = json_decode($cake['detected'], true);

		$cakeComments = implode("\n", $cakeCommentsList);

		$pl = json_decode($cake['processes'], true);

		echo '<p align="center"><font size="5"><i class="fa fa-birthday-cake"></i>	Edit cake#'.$id.'</font></p>';

		echo '<div class="text-center">
			<a href="index.php?p=128&uid='.$cake['userid'].'" type="button" class="btn btn-primary col-md-2">Go back</a>
			<a href="/u/'.$cake['userid'].'" type="button" class="btn btn-info col-md-8">Profile</a>
			<a href="index.php?p=103&id='.$cake['userid'].'" type="button" class="btn btn-warning col-md-2">Edit User</a>
		</div>';

		echo '<table class="table table-striped table-hover table-center"><tbody>';

		echo '<form id="system-settings-form" action="submit.php" method="POST"></form>';

		echo '<tr>
		<td width=1>ID</td>
		<td><p class="text-center"><input type="number" name="id" class="form-control" value="'.$id.'" readonly=""></p></td>
		</tr>';

		echo '<tr>
		<td>Username</td>
		<td><p class="text-center"><input type="text" name="username" class="form-control" value="'.$user['username'].'" readonly=""></p></td>
		</tr>';

		echo '<tr>
		<td>Score ID</td>
		<td><p class="text-center"><input type="number" name="scoreid" class="form-control" value="'.$cake['score_id'].'" readonly=""></p>
		<div class="text-center">
			<a href="'.Fringuellina::getBeatmapUrl($beatmap['beatmap_id']).'" type="button" class="btn btn-success">Download Beatmap</a>
			<a href="http://'.Fringuellina::getMainDomain().'/web/replays/'.$cake['score_id'].'" type="button" class="btn btn-primary">Download Replay</a>
		</div>
		</tr>';

		echo '<tr>
		<td>Flags</td>
		<td>
		<p class="text-center"><input type="number" name="flags" class="form-control" value="'.$cake['flags'].'"></p>
		<p class="text-center"><input type="text" name="flags_string" class="form-control" value="'.$flags.'" readonly=""></p>
		</td>
		</tr>';

		echo '<tr>
		<td>Cake Comments</td>
		<td>
		<textarea name="cake_comments" class="form-control" style="overflow: auto; resize: vertical; height: 209px; margin-top: 0px; margin-bottom: 0px;">'.$cakeComments.'</textarea>
		</td>
		</tr>';
		
		echo '<tr>
		<td>Cake Ingredients</td>
		<td>';

		echo '<div class="text-center">
		<a class="btn btn-success" id="collapseVisualB">Show Visual</a>
		<a class="btn btn-success" id="collapseJsonPB">Show Json Prettified</a>
		<a class="btn btn-success" id="collapseJsonRB">Show Json Raw</a>';

		echo '<div class="collapse" id="collapseVisual">';
		foreach ($pl as $item){
			if ($item['hash'] == null && $item['path'] == null && $item['title'] == null)
				if (in_array($item['file'], ['svchost', 'SearchIndexer', 'chrome', 'smss', 'SearchUI', 'csrss', 'RuntimeBroker', 'spoolsv', 'SettingSyncHost', 'Memory Compression', 'conhost', 'lsass', 'conhost', 'dwm', 'rundll32', 'dllhost', 'Idle']))
					continue;

			$c = "primary";
			//Do some check to see if it is type WARNING (flagged)
			foreach ($eggs as $egg){
				if ($egg["is_regex"]){
					if (preg_match('/'.$egg["value"].'/', $item[$egg["type"]])){
						if ($egg["ban"])
							$c = "danger";
						else
							$c = "warning";
					}		
				}
				else if ($egg["value"] == $item[$egg["type"]]){
					if ($egg["ban"])
						$c = "danger";
					else
						$c = "warning";
				}
			}

			echo '<a class="btn btn-block btn-'.$c.'">'.$item['file'].'<br>'.$item['hash'].'<br>'.$item['path'].'<br>'.$item['title'].'<br></a>';
		}
		echo '</div>
		<div class="collapse" id="collapseJsonP">
			<textarea name="jsonP" class="form-control" style="overflow: auto; resize: vertical; height: 800px; margin-top: 0px; margin-bottom: 0px;" readonly="">'.json_encode($pl, JSON_PRETTY_PRINT).'</textarea>
		</div>
		<div class="collapse" id="collapseJsonR">
			<textarea name="jsonR" class="form-control" style="overflow: auto; resize: vertical; height: 800px; margin-top: 0px; margin-bottom: 0px;" readonly="">'.json_encode($pl).'</textarea>
		</div>
		</div>';

		echo '</td>
		</tr>';

		echo '</tbody></table>';

		echo '</div></div>';

		echo '<script>
		document.body.onload = function(){
			$("#collapseVisualB").click(function(){
				$("#collapseVisual").collapse("toggle");
				$("#collapseJsonP").collapse("hide");
				$("#collapseJsonR").collapse("hide");
			});

			$("#collapseJsonPB").click(function(){
				$("#collapseVisual").collapse("hide");
				$("#collapseJsonP").collapse("toggle");
				$("#collapseJsonR").collapse("hide");
			});

			$("#collapseJsonRB").click(function(){
				$("#collapseVisual").collapse("hide");
				$("#collapseJsonP").collapse("hide");
				$("#collapseJsonR").collapse("toggle");
			});
		};
		</script>';
	}

    public static function PrintCakesSummary(){
		// Print stuff
		echo '<div id="wrapper">';
		printAdminSidebar();
		echo '<div id="page-content-wrapper">';
		// Maintenance check
		P::MaintenanceStuff();
		// Print Success if set
		if (isset($_GET['s']) && !empty($_GET['s'])) {
			P::SuccessMessageStaccah($_GET['s']);
		}
		// Print Exception if set
		if (isset($_GET['e']) && !empty($_GET['e'])) {
			P::ExceptionMessageStaccah($_GET['e']);
        }
    }

    public static function PrintEditCake(){

    }

    public static function RAPButton(){
		echo '<li><a href="index.php?p=128"><i class="fa fa-birthday-cake"></i>	Cakes</a></li>';
    }

    public static function RAPCakesListButton(){
		echo '<li><a href="index.php?p=130"><i class="fa fa-book"></i>	Cake recipes</a></li>';
    }

	//TODO; Make seperate php file to redirect to correct page with custom lookups etc.
	public static function PrintLookupUserModule(){
		echo '<div class="modal fade" id="quickLookupUserModal" tabindex="-1" role="dialog" aria-labelledby="quickLookupUserModal">
		<div class="modal-dialog">
		<div class="modal-content">
		<div class="modal-header">
		<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		<h4 class="modal-title" id="quickLookupUserModal">Quick lookup user</h4>
		</div>
		<div class="modal-body">
		<p>
		<form id="quick-lookup-user-id" action="submit.php" method="POST">
		<input name="action" value="toggleCake" hidden>
		<input name="extra_action" value="quickLookupUser" hidden>
		<div class="input-group">
		<span class="input-group-addon" id="basic-addon1"><span class="glyphicon glyphicon-user" aria-hidden="true"></span></span>
		<input type="text" name="data" class="form-control" placeholder="Username" aria-describedby="basic-addon1" required>
		</div>
		</form>
		</p>
		</div>
		<div class="modal-footer">
		<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
		<button type="submit" form="quick-lookup-user-id" class="btn btn-primary">Check</button>
		</div>
		</div>
		</div>
		</div>';
	}

	public static function PrintLookupCakeModule(){
		echo '<div class="modal fade" id="quickLookupScoreidModal" tabindex="-1" role="dialog" aria-labelledby="quickLookupScoreidModal">
		<div class="modal-dialog">
		<div class="modal-content">
		<div class="modal-header">
		<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		<h4 class="modal-title" id="quickLookupScoreidModal">Quick lookup cake</h4>
		</div>
		<div class="modal-body">
		<p>
		<form id="quick-lookup-score-id" action="submit.php" method="POST">
		<input name="action" value="toggleCake" hidden>
		<input name="extra_action" value="quickLookupScoreid" hidden>
		<div class="input-group">
		<span class="input-group-addon" id="basic-addon1"><span class="glyphicon glyphicon-tag" aria-hidden="true"></span></span>
		<input type="text" onkeypress="return event.charCode >= 48 && event.charCode <= 57" name="data" class="form-control" placeholder="Score ID" aria-describedby="basic-addon1" required>
		</div>
		</form>
		</p>
		</div>
		<div class="modal-footer">
		<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
		<button type="submit" form="quick-lookup-score-id" class="btn btn-primary">Check</button>
		</div>
		</div>
		</div>
		</div>';
	}

	public static function getMainDomain(){
		$host_names = explode(".", $_SERVER['SERVER_NAME']);
		unset($host_names[0]);
		$main_domain = implode(".", $host_names);

		return $main_domain;
	}

	public static function getBeatmapUrl($id){
		return "http://".Fringuellina::getMainDomain()."/b/".$id;
	}

	const IGNORE_FLAGS = BadFlags::CLEAN | BadFlags::INCORRECT_MOD;
	public static function makeFlagString($i){
		$flags = [];

		$ref = new ReflectionClass("BadFlags");
		$arr = $ref->getConstants();

		foreach ($arr as $flag)
		{
			if (($i & $flag) != 0 && ($i & ~Fringuellina::IGNORE_FLAGS))
				array_push($flags, array_search($flag, $arr));
		}

		unset($ref);

		return implode(" | ", $flags);
	}

	//submit.php
	//toggleCake
	//I am using this for some extended stuff
	public static function toggleCake(){
		if (isset($_POST['extra_action']) && !empty($_POST['extra_action'])) {
			$action = $_POST['extra_action'];
		} elseif (isset($_GET['extra_action']) && !empty($_GET['extra_action'])) {
			$action = $_GET['extra_action'];
		} else {
			Fringuellina::actualToggleCake();
			die();
		}

		try{
			if (isset($_POST['data']) && !empty($_POST['data'])) {
				$data = $_POST['data'];
			} elseif (isset($_GET['data']) && !empty($_GET['data'])) {
				$data = $_GET['data'];
			} else {
				throw new Exception("Couldn't find data parameter");
			}

			switch($action){
				case "quickLookupUser":
					Fringuellina::QuickLookupUser();
					break;
				case "quickLookupScoreid":
					Fringuellina::QuickLookupScoreid();
					break;
			}
		}
		catch(Exception $e) {
			// Redirect to Exception page
			redirect('index.php?p=99&e='.$e->getMessage());
		}
	}

	public static function actualToggleCake(){

	}

	public static function QuickLookupUser(){
		$uid = current($GLOBALS['db']->fetch('SELECT id FROM users WHERE username = ?', [$_POST['data']]));
		header('Location: /index.php?p=128&uid='.$uid);
		exit();
	}

	public static function QuickLookupScoreid(){
		$id = current($GLOBALS['db']->fetch('SELECT id FROM cakes WHERE score_id = ?', [$_POST['data']]));
		header('Location: /index.php?p=129&id='.$id);
		exit();
	}

	//removeCake
	public static function RemoveCake(){

    }
	//saveCake
    public static function EditCake(){

	}
}

class BadFlags{
	const CLEAN = 0;
	const SPEED = 1 << 1;
	const INCORRECT_MOD = 1 << 2;
	const MULTIPLE_OSU_CLIENTS = 1 << 3;
	const CHECKSUM_FAIL = 1 << 4;
	const FLASHLIGHT_CHECKSUM_FAIL = 1 << 5;
	const OSU_CHECKSUM = 1 << 6;
	const MISSING_PL = 1 << 7;
	const FLASHLIGHT_IMAGE = 1 << 8;
	const SPINNER = 1 << 9;
	const TRANSPARENT_WINDOW = 1 << 10;
	const FAST_PRESS = 1 << 11;
	const RAW_MOUSE_DISCREPANCY = 1 << 12;
	const RAW_KEYBOARD_DISCREPANCY = 1 << 13;
}
?>