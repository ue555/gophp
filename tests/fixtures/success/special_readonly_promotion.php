<?php
readonly class LinkedNode {
    public function __construct(public ?self $next = null) {}
}
$last = new LinkedNode;
$first = new LinkedNode($last);
echo $first->next === $last ? 'linked' : 'wrong';
