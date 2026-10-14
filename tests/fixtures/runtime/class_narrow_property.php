<?php
echo "before\n";
class VisiblePropertyParent { protected $value; }
class VisiblePropertyChild extends VisiblePropertyParent { private $value; }
