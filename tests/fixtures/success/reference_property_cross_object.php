<?php
class RefCross { public $left = 1; public $right = 2; }
$first = new RefCross;
$second = new RefCross;
$first->left =& $second->right;
$first->left = 7;
echo $first->left, '|', $second->right;
