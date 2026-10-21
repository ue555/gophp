<?php
trait AliasSource { public function original() { return 1; } }
class AliasPriority {
    use AliasSource { original as selected; }
    public function selected() { return 2; }
}
$x = new AliasPriority;
echo $x->original(), '|', $x->selected(), "\n";
