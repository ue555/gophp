<?php
interface RootTask { public function root(); }
interface LeftTask extends RootTask { public function left(); }
interface RightTask { public function right(); }
interface JointTask extends LeftTask, RightTask {}
class TaskBase implements JointTask {
    public function root() { return 1; }
    public function left() { return 2; }
    public function right() { return 3; }
}
class TaskChild extends TaskBase {}
$x = new TaskChild;
echo $x->root(), $x->left(), $x->right(), '|';
echo $x instanceof RootTask, $x instanceof LeftTask, $x instanceof RightTask, $x instanceof JointTask, "\n";
