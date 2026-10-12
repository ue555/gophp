<?php
class Returns {
    function early($x) {
        if ($x) { if (true) return 'yes'; echo 'bad'; }
        elseif (false) return 'bad';
        else return;
        echo 'bad';
    }
    function implicit() { $x = 1; }
}
$o = new Returns;
echo $o->early(true), ':', ($o->early(false) === null) + 0, ($o->implicit() === null) + 0, "\n";
