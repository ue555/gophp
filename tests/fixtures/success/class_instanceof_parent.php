<?php
class RootType {}
class ChildType extends RootType {}
$item = new ChildType;
echo ($item instanceof ChildType) ? 'child' : 'no';
echo '|', ($item instanceof RootType) ? 'root' : 'no';
echo '|', (null instanceof RootType) ? 'bad' : 'null', "\n";
