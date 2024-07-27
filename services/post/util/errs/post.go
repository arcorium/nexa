package errs

import "errors"

var ErrPostWithNoVersion = errors.New("post malformed, it has no versions")

var ErrGetBlockedUserPost = errors.New("you are not allowed to get post from this user")

var ErrNoRightAccessToGetPost = errors.New("you are not allowed to get this post")
