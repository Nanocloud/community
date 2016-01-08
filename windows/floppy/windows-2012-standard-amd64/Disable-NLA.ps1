# Poweshell script to disable NLA option for our collection
Set-RDSessionCollectionConfiguration -CollectionName collection -AuthenticateUsingNLA 0
