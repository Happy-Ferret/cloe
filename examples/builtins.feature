Feature: Built-in functions
  Scenario: Get types of values
    Given a file named "main.coel" with:
    """
    (seq
      (write (typeOf true))
      (write (typeOf {"key" "value"}))
      (write (typeOf []))
      (write (typeOf nil))
      (write (typeOf 42))
      (write (typeOf "foo"))
      (write (typeOf +))
      (write (typeOf (partial + 1))))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    bool
    dict
    list
    nil
    number
    string
    function
    function
    """

  Scenario: Map a function to a list
    Given a file named "main.coel" with:
    """
    (write (map (\ (x) (* x x)) [1 2 3]))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [1 4 9]
    """

  Scenario: Calculate indices of elements in a list
    Given a file named "main.coel" with:
    """
    (let l [1 2 3 42 -3 "foo"])
    (seq
      (write (indexOf l 42))
      (write (indexOf l 2))
      (write (indexOf l "foo")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    3
    1
    5
    """

  Scenario: Use multiple conditions with if function
    Given a file named "main.coel" with:
    """
    (def (no) (write "No"))

    (if false no true (write "Yes") false no no)
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Yes
    """

  Scenario: Slice lists
    Given a file named "main.coel" with:
    """
    (seq
      (write (slice [1 2 3]))
      (write (slice [1 2 3] 0))
      (write (slice [1 2 3] 1 3))
      (write (slice [1 2 3] 0 2))
      (write (slice [1 2 3] 1))
      (write (slice [1 2 3] . start 1))
      (write (slice [1 2 3] . end 1))
      (write (slice [1 2 3] . start 2))
      (write (slice [1 2 3] . start 3))
      (write (slice [1 2 3] . start 4))
      (write (slice [1 2 3] . end 0)))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [1 2 3]
    [1 2 3]
    [2 3]
    [1 2]
    [2 3]
    [2 3]
    [1]
    [3]
    []
    []
    []
    """
