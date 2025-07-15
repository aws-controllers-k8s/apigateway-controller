    if err := updateDomainNameInput(desired, latest, input, delta); err != nil {
        return nil, err
    }