-- rollback for 000005 (drop avatar storage policies + bucket; profile rows keep avatar_path text)
drop policy if exists "avatars owner delete" on storage.objects;
drop policy if exists "avatars owner update" on storage.objects;
drop policy if exists "avatars owner insert" on storage.objects;
drop policy if exists "avatars public read" on storage.objects;
delete from storage.buckets where id = 'avatars';
