<?php
$inf = 1e308 * 10;
$nan = $inf - $inf;
$t = true; $f = false; $nil = null;
// $inf versus 1e308: == != < <= > >=
echo (($inf) == (1e308)) + 0, (($inf) != (1e308)) + 0, (($inf) < (1e308)) + 0, (($inf) <= (1e308)) + 0, (($inf) > (1e308)) + 0, (($inf) >= (1e308)) + 0, "\n";
// $inf versus $inf: == != < <= > >=
echo (($inf) == ($inf)) + 0, (($inf) != ($inf)) + 0, (($inf) < ($inf)) + 0, (($inf) <= ($inf)) + 0, (($inf) > ($inf)) + 0, (($inf) >= ($inf)) + 0, "\n";
// -$inf versus $inf: == != < <= > >=
echo ((-$inf) == ($inf)) + 0, ((-$inf) != ($inf)) + 0, ((-$inf) < ($inf)) + 0, ((-$inf) <= ($inf)) + 0, ((-$inf) > ($inf)) + 0, ((-$inf) >= ($inf)) + 0, "\n";
// $inf versus 'INF': == != < <= > >=
echo (($inf) == ('INF')) + 0, (($inf) != ('INF')) + 0, (($inf) < ('INF')) + 0, (($inf) <= ('INF')) + 0, (($inf) > ('INF')) + 0, (($inf) >= ('INF')) + 0, "\n";
// -$inf versus '-INF': == != < <= > >=
echo ((-$inf) == ('-INF')) + 0, ((-$inf) != ('-INF')) + 0, ((-$inf) < ('-INF')) + 0, ((-$inf) <= ('-INF')) + 0, ((-$inf) > ('-INF')) + 0, ((-$inf) >= ('-INF')) + 0, "\n";
// $inf versus '1e999': == != < <= > >=
echo (($inf) == ('1e999')) + 0, (($inf) != ('1e999')) + 0, (($inf) < ('1e999')) + 0, (($inf) <= ('1e999')) + 0, (($inf) > ('1e999')) + 0, (($inf) >= ('1e999')) + 0, "\n";
// $nan versus $nan: == != < <= > >=
echo (($nan) == ($nan)) + 0, (($nan) != ($nan)) + 0, (($nan) < ($nan)) + 0, (($nan) <= ($nan)) + 0, (($nan) > ($nan)) + 0, (($nan) >= ($nan)) + 0, "\n";
// $nan versus 0: == != < <= > >=
echo (($nan) == (0)) + 0, (($nan) != (0)) + 0, (($nan) < (0)) + 0, (($nan) <= (0)) + 0, (($nan) > (0)) + 0, (($nan) >= (0)) + 0, "\n";
// 0 versus $nan: == != < <= > >=
echo ((0) == ($nan)) + 0, ((0) != ($nan)) + 0, ((0) < ($nan)) + 0, ((0) <= ($nan)) + 0, ((0) > ($nan)) + 0, ((0) >= ($nan)) + 0, "\n";
// $nan versus $inf: == != < <= > >=
echo (($nan) == ($inf)) + 0, (($nan) != ($inf)) + 0, (($nan) < ($inf)) + 0, (($nan) <= ($inf)) + 0, (($nan) > ($inf)) + 0, (($nan) >= ($inf)) + 0, "\n";
// $inf versus $nan: == != < <= > >=
echo (($inf) == ($nan)) + 0, (($inf) != ($nan)) + 0, (($inf) < ($nan)) + 0, (($inf) <= ($nan)) + 0, (($inf) > ($nan)) + 0, (($inf) >= ($nan)) + 0, "\n";
// $nan versus 'NAN': == != < <= > >=
echo (($nan) == ('NAN')) + 0, (($nan) != ('NAN')) + 0, (($nan) < ('NAN')) + 0, (($nan) <= ('NAN')) + 0, (($nan) > ('NAN')) + 0, (($nan) >= ('NAN')) + 0, "\n";
// 'NAN' versus $nan: == != < <= > >=
echo (('NAN') == ($nan)) + 0, (('NAN') != ($nan)) + 0, (('NAN') < ($nan)) + 0, (('NAN') <= ($nan)) + 0, (('NAN') > ($nan)) + 0, (('NAN') >= ($nan)) + 0, "\n";
// $nan versus '0': == != < <= > >=
echo (($nan) == ('0')) + 0, (($nan) != ('0')) + 0, (($nan) < ('0')) + 0, (($nan) <= ('0')) + 0, (($nan) > ('0')) + 0, (($nan) >= ('0')) + 0, "\n";
// '' versus $nan: == != < <= > >=
echo (('') == ($nan)) + 0, (('') != ($nan)) + 0, (('') < ($nan)) + 0, (('') <= ($nan)) + 0, (('') > ($nan)) + 0, (('') >= ($nan)) + 0, "\n";
// $nan versus true: == != < <= > >=
echo (($nan) == (true)) + 0, (($nan) != (true)) + 0, (($nan) < (true)) + 0, (($nan) <= (true)) + 0, (($nan) > (true)) + 0, (($nan) >= (true)) + 0, "\n";
// true versus $nan: == != < <= > >=
echo ((true) == ($nan)) + 0, ((true) != ($nan)) + 0, ((true) < ($nan)) + 0, ((true) <= ($nan)) + 0, ((true) > ($nan)) + 0, ((true) >= ($nan)) + 0, "\n";
// $nan versus $t: == != < <= > >=
echo (($nan) == ($t)) + 0, (($nan) != ($t)) + 0, (($nan) < ($t)) + 0, (($nan) <= ($t)) + 0, (($nan) > ($t)) + 0, (($nan) >= ($t)) + 0, "\n";
// $t versus $nan: == != < <= > >=
echo (($t) == ($nan)) + 0, (($t) != ($nan)) + 0, (($t) < ($nan)) + 0, (($t) <= ($nan)) + 0, (($t) > ($nan)) + 0, (($t) >= ($nan)) + 0, "\n";
// $nan versus false: == != < <= > >=
echo (($nan) == (false)) + 0, (($nan) != (false)) + 0, (($nan) < (false)) + 0, (($nan) <= (false)) + 0, (($nan) > (false)) + 0, (($nan) >= (false)) + 0, "\n";
// false versus $nan: == != < <= > >=
echo ((false) == ($nan)) + 0, ((false) != ($nan)) + 0, ((false) < ($nan)) + 0, ((false) <= ($nan)) + 0, ((false) > ($nan)) + 0, ((false) >= ($nan)) + 0, "\n";
// $nan versus $f: == != < <= > >=
echo (($nan) == ($f)) + 0, (($nan) != ($f)) + 0, (($nan) < ($f)) + 0, (($nan) <= ($f)) + 0, (($nan) > ($f)) + 0, (($nan) >= ($f)) + 0, "\n";
// $f versus $nan: == != < <= > >=
echo (($f) == ($nan)) + 0, (($f) != ($nan)) + 0, (($f) < ($nan)) + 0, (($f) <= ($nan)) + 0, (($f) > ($nan)) + 0, (($f) >= ($nan)) + 0, "\n";
// $nan versus null: == != < <= > >=
echo (($nan) == (null)) + 0, (($nan) != (null)) + 0, (($nan) < (null)) + 0, (($nan) <= (null)) + 0, (($nan) > (null)) + 0, (($nan) >= (null)) + 0, "\n";
// null versus $nan: == != < <= > >=
echo ((null) == ($nan)) + 0, ((null) != ($nan)) + 0, ((null) < ($nan)) + 0, ((null) <= ($nan)) + 0, ((null) > ($nan)) + 0, ((null) >= ($nan)) + 0, "\n";
// $nan versus $nil: == != < <= > >=
echo (($nan) == ($nil)) + 0, (($nan) != ($nil)) + 0, (($nan) < ($nil)) + 0, (($nan) <= ($nil)) + 0, (($nan) > ($nil)) + 0, (($nan) >= ($nil)) + 0, "\n";
// $nil versus $nan: == != < <= > >=
echo (($nil) == ($nan)) + 0, (($nil) != ($nan)) + 0, (($nil) < ($nan)) + 0, (($nil) <= ($nan)) + 0, (($nil) > ($nan)) + 0, (($nil) >= ($nan)) + 0, "\n";
