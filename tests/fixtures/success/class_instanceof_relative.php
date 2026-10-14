<?php
class RelativeRoot {
    public function kinds($x) { echo $x instanceof self, '|', $x instanceof static; }
}
class RelativeChild extends RelativeRoot {
    public function parentCheck($x) { return $x instanceof parent; }
}
$child = new RelativeChild;
$child->kinds(new RelativeRoot);
echo '|', $child->parentCheck($child), "\n";
