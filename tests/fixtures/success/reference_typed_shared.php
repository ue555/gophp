<?php
class TypedShared {
    public int $left = 1;
    public int $right = 2;
    public $plain = 0;
}
$box = new TypedShared;
$alias =& $box->left;
$box->right =& $alias;
$box->plain =& $alias;
$box->plain = '8';
echo $box->left, '|', $box->right, '|', gettype($box->plain);
